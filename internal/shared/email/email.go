package email

import (
	"fmt"
	"html"
	"strings"
)

const SenderDisplayName = "Macabi Madrijim"

func Layout(headerColor, title, body string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="es">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"></head>
<body style="margin:0;padding:0;background:#f4f6f9;font-family:Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f6f9;padding:32px 0;">
    <tr><td align="center">
      <table width="600" cellpadding="0" cellspacing="0" style="max-width:600px;width:100%%;background:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
        <tr>
          <td style="background:%s;padding:28px 32px;">
            <p style="margin:0;font-size:13px;color:rgba(255,255,255,0.8);letter-spacing:1px;text-transform:uppercase;">Macabi Madrijim</p>
            <h1 style="margin:6px 0 0;font-size:22px;color:#ffffff;font-weight:700;">%s</h1>
          </td>
        </tr>
        <tr>
          <td style="padding:32px;">
            %s
          </td>
        </tr>
        <tr>
          <td style="background:#f4f6f9;padding:20px 32px;border-top:1px solid #e5e7eb;">
            <p style="margin:0;font-size:12px;color:#9ca3af;text-align:center;">Este mensaje fue generado automáticamente por Macabi Madrijim. Por favor no respondas este correo.</p>
          </td>
        </tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`, headerColor, html.EscapeString(title), body)
}

type DetailRow struct {
	Label string
	Value string
}

func DetailsCard(rows []DetailRow) string {
	var b strings.Builder
	b.WriteString(`<table width="100%" cellpadding="0" cellspacing="0" style="background:#f8fafc;border:1px solid #e5e7eb;border-radius:6px;margin:20px 0;">`)
	for i, r := range rows {
		border := "border-bottom:1px solid #e5e7eb;"
		if i == len(rows)-1 {
			border = ""
		}
		b.WriteString(fmt.Sprintf(
			`<tr><td style="padding:14px 20px;%s"><span style="font-size:12px;color:#6b7280;text-transform:uppercase;letter-spacing:0.5px;">%s</span><br><span style="font-size:16px;color:#111827;font-weight:600;">%s</span></td></tr>`,
			border, html.EscapeString(r.Label), html.EscapeString(r.Value),
		))
	}
	b.WriteString(`</table>`)
	return b.String()
}

func CTAButton(url, label, color string) string {
	return fmt.Sprintf(`
<table width="100%%" cellpadding="0" cellspacing="0" style="margin-top:24px;">
  <tr>
    <td align="center">
      <a href="%s" style="display:inline-block;background:%s;color:#ffffff;font-size:15px;font-weight:700;text-decoration:none;padding:14px 40px;border-radius:6px;">%s</a>
    </td>
  </tr>
  <tr>
    <td align="center" style="padding-top:12px;">
      <span style="font-size:12px;color:#9ca3af;">O copiá este enlace: <a href="%s" style="color:#6b7280;">%s</a></span>
    </td>
  </tr>
</table>`, url, color, html.EscapeString(label), url, url)
}
