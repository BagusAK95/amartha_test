package mail

import "context"

type IMailUsecase interface {
	Send(ctx context.Context, req MailSendRequest)
}
