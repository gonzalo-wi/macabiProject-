package newsports

import "context"

type NewsMailer interface {
	NotifyMembersNewNews(ctx context.Context, recipientEmails []string, title, summary, newsID string) error
}
