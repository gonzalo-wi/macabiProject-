package expensesstorage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type SupabaseSigner struct {
	BaseURL string
	APIKey  string
	Bucket  string
	Client  *http.Client
}

func (s *SupabaseSigner) client() *http.Client {
	if s.Client != nil {
		return s.Client
	}
	return http.DefaultClient
}

func encodeObjectPath(bucket, key string) string {
	parts := append([]string{bucket}, strings.Split(key, "/")...)
	out := make([]string, len(parts))
	for i, p := range parts {
		out[i] = url.PathEscape(p)
	}
	return strings.Join(out, "/")
}

func (s *SupabaseSigner) resolveStorageURL(relativeOrAbsolute string) string {
	relativeOrAbsolute = strings.TrimSpace(relativeOrAbsolute)
	if relativeOrAbsolute == "" {
		return ""
	}
	if strings.HasPrefix(relativeOrAbsolute, "http://") || strings.HasPrefix(relativeOrAbsolute, "https://") {
		return relativeOrAbsolute
	}
	base := strings.TrimRight(s.BaseURL, "/")
	if !strings.HasPrefix(relativeOrAbsolute, "/") {
		relativeOrAbsolute = "/" + relativeOrAbsolute
	}
	if strings.HasPrefix(relativeOrAbsolute, "/storage/v1/") {
		return base + relativeOrAbsolute
	}
	if strings.HasPrefix(relativeOrAbsolute, "/object/") {
		return base + "/storage/v1" + relativeOrAbsolute
	}
	return base + relativeOrAbsolute
}

func (s *SupabaseSigner) CreateSignedUploadURL(ctx context.Context, objectKey, contentType string) (string, error) {
	_ = contentType // callers set Content-Type on PUT; signer path kept for symmetry with interface
	path := encodeObjectPath(s.Bucket, objectKey)
	u := strings.TrimRight(s.BaseURL, "/") + "/storage/v1/object/upload/sign/" + path

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(nil))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("apikey", s.APIKey)
	resp, err := s.client().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("supabase upload sign failed: status %d body %s", resp.StatusCode, string(body))
	}

	var decoded struct {
		U string `json:"URL"`
		L string `json:"url"`
	}
	if err := json.Unmarshal(body, &decoded); err != nil {
		return "", err
	}
	relative := decoded.U
	if relative == "" {
		relative = decoded.L
	}
	if relative == "" {
		return "", fmt.Errorf("supabase upload sign: empty url")
	}

	return s.resolveStorageURL(relative), nil
}

func (s *SupabaseSigner) CreateSignedDownloadURL(ctx context.Context, objectKey string, expiresSec int) (string, error) {
	path := encodeObjectPath(s.Bucket, objectKey)
	u := strings.TrimRight(s.BaseURL, "/") + "/storage/v1/object/sign/" + path

	payload := fmt.Sprintf(`{"expiresIn":%d}`, expiresSec)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("apikey", s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("supabase download sign failed: status %d body %s", resp.StatusCode, string(body))
	}

	var decoded struct {
		A string `json:"signedURL"`
		B string `json:"signedUrl"`
	}
	if err := json.Unmarshal(body, &decoded); err != nil {
		return "", err
	}
	out := decoded.A
	if out == "" {
		out = decoded.B
	}
	return s.resolveStorageURL(out), nil
}

func (s *SupabaseSigner) UploadObject(ctx context.Context, objectKey, contentType string, body io.Reader) error {
	path := encodeObjectPath(s.Bucket, objectKey)
	u := strings.TrimRight(s.BaseURL, "/") + "/storage/v1/object/" + path

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("apikey", s.APIKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "true")

	resp, err := s.client().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("supabase object upload failed: status %d body %s", resp.StatusCode, string(respBody))
	}
	return nil
}

func (s *SupabaseSigner) DeleteObject(ctx context.Context, objectKey string) error {
	path := encodeObjectPath(s.Bucket, objectKey)
	u := strings.TrimRight(s.BaseURL, "/") + "/storage/v1/object/" + path

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("apikey", s.APIKey)

	resp, err := s.client().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("supabase object delete failed: status %d body %s", resp.StatusCode, string(respBody))
	}
	return nil
}
