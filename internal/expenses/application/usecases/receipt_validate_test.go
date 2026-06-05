package expensesusecases

import (
	"bytes"
	"errors"
	"testing"

	expensesdomain "macabi-back/internal/expenses/domain"
)

func TestValidateReceiptBytes_rejectsOversize(t *testing.T) {
	body := bytes.Repeat([]byte{0xff, 0xd8, 0xff}, MaxReceiptBytes/3+1)
	_, _, err := ValidateReceiptBytes(bytes.NewReader(body), int64(len(body)), "image/jpeg")
	if !errors.Is(err, expensesdomain.ErrReceiptTooLarge) {
		t.Fatalf("want ErrReceiptTooLarge, got %v", err)
	}
}

func TestValidateReceiptBytes_rejectsExecutableDisguisedAsJPEG(t *testing.T) {
	body := []byte("MZ\x90\x00fake exe content")
	_, _, err := ValidateReceiptBytes(bytes.NewReader(body), int64(len(body)), "image/jpeg")
	if !errors.Is(err, expensesdomain.ErrInvalidMimeType) {
		t.Fatalf("want ErrInvalidMimeType, got %v", err)
	}
}

func TestValidateReceiptBytes_acceptsPNG(t *testing.T) {
	body := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0, 0, 0, 0x0d}
	data, ct, err := ValidateReceiptBytes(bytes.NewReader(body), int64(len(body)), "image/png")
	if err != nil {
		t.Fatal(err)
	}
	if ct != "image/png" {
		t.Fatalf("ct=%s", ct)
	}
	if len(data) != len(body) {
		t.Fatalf("len=%d", len(data))
	}
}

func TestValidateReceiptBytes_acceptsPDF(t *testing.T) {
	body := []byte("%PDF-1.4\n%fake")
	_, ct, err := ValidateReceiptBytes(bytes.NewReader(body), int64(len(body)), "application/pdf")
	if err != nil {
		t.Fatal(err)
	}
	if ct != "application/pdf" {
		t.Fatalf("ct=%s", ct)
	}
}

func TestValidateReceiptBytes_rejectsMismatchedDeclaredType(t *testing.T) {
	body := []byte("%PDF-1.4\n")
	_, _, err := ValidateReceiptBytes(bytes.NewReader(body), int64(len(body)), "image/png")
	if !errors.Is(err, expensesdomain.ErrInvalidMimeType) {
		t.Fatalf("want ErrInvalidMimeType, got %v", err)
	}
}
