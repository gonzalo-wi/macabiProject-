package expensesusecases

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	expensesdomain "macabi-back/internal/expenses/domain"
)

// ValidateReceiptBytes reads up to MaxReceiptBytes+1, detects real MIME from magic bytes,
// and rejects spoofed Content-Type headers. Returns canonical content-type for storage.
func ValidateReceiptBytes(body io.Reader, declaredSize int64, declaredContentType string) ([]byte, string, error) {
	if body == nil {
		return nil, "", expensesdomain.ErrMissingRequiredField
	}
	if declaredSize < 0 {
		return nil, "", expensesdomain.ErrReceiptTooLarge
	}
	if declaredSize > MaxReceiptBytes {
		return nil, "", expensesdomain.ErrReceiptTooLarge
	}

	limited := io.LimitReader(body, MaxReceiptBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, "", expensesdomain.ErrReceiptAttachFailed
	}
	if len(data) == 0 {
		return nil, "", expensesdomain.ErrMissingRequiredField
	}
	if len(data) > MaxReceiptBytes {
		return nil, "", expensesdomain.ErrReceiptTooLarge
	}

	detected := detectReceiptMIME(data)
	if detected == "" {
		return nil, "", expensesdomain.ErrInvalidMimeType
	}

	declared := normalizeDeclaredMIME(declaredContentType)
	if declared != "" && declared != detected {
		// Allow image/jpeg vs image/jpg style drift only when both map to same suffix
		if receiptSuffix(declared) != receiptSuffix(detected) {
			return nil, "", expensesdomain.ErrInvalidMimeType
		}
	}

	return data, detected, nil
}

func normalizeDeclaredMIME(ct string) string {
	ct = strings.TrimSpace(strings.ToLower(ct))
	if ct == "" {
		return ""
	}
	if i := strings.Index(ct, ";"); i >= 0 {
		ct = strings.TrimSpace(ct[:i])
	}
	return ct
}

func receiptSuffix(mime string) string {
	s, ok := ReceiptMIME[mime]
	if !ok {
		return ""
	}
	return s
}

// detectReceiptMIME identifies allowed types from file header (not client Content-Type).
func detectReceiptMIME(data []byte) string {
	if len(data) >= 5 && string(data[:5]) == "%PDF-" {
		return "application/pdf"
	}
	if len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
		return "image/webp"
	}
	if len(data) >= 8 && data[0] == 0x89 && string(data[1:4]) == "PNG" {
		return "image/png"
	}
	if len(data) >= 3 && data[0] == 0xff && data[1] == 0xd8 && data[2] == 0xff {
		return "image/jpeg"
	}

	sniff := http.DetectContentType(data)
	sniff = normalizeDeclaredMIME(sniff)
	if _, ok := ReceiptMIME[sniff]; ok {
		return sniff
	}
	return ""
}

// bytesReader is a tiny helper for tests and upload after validation.
func bytesReader(data []byte) *bytes.Reader {
	return bytes.NewReader(data)
}
