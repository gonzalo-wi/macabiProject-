package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const brevoAPIURL = "https://api.brevo.com/v3/smtp/email"

type BrevoRecipient struct {
	Email string `json:"email"`
}

type brevoSender struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type brevoRequest struct {
	Sender      brevoSender      `json:"sender"`
	To          []BrevoRecipient `json:"to"`
	Subject     string           `json:"subject"`
	HTMLContent string           `json:"htmlContent"`
}

func BrevoSend(ctx context.Context, client *http.Client, apiKey, from string, to []BrevoRecipient, subject, html string) error {
	if apiKey == "" || from == "" {
		return nil
	}
	body := brevoRequest{
		Sender:      brevoSender{Email: from, Name: SenderDisplayName},
		To:          to,
		Subject:     subject,
		HTMLContent: html,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal brevo payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, brevoAPIURL, bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("brevo request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("brevo http: %w", err)
	}
	defer res.Body.Close()
	respBody, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("brevo: status %d: %s", res.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return nil
}
