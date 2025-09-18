package mail

import "context"

type IMailUsecase interface {
	Send(ctx context.Context, req MailSendRequest)
}

// Move to dto folder
type MailSendRequest struct {
	To       string
	Subject  string
	Template string
	Data     map[string]any
}
