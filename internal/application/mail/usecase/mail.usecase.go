package usecase

import (
	"context"
	"log"

	"github.com/BagusAK95/amarta_test/internal/domain/mail"
	mailsender "github.com/BagusAK95/amarta_test/internal/infrastructure/mail"
	"github.com/opentracing/opentracing-go"
)

type mailUsecase struct {
	mailSender *mailsender.Sender
}

func NewMailUsecase(mailSender *mailsender.Sender) mail.IMailUsecase {
	return &mailUsecase{
		mailSender: mailSender,
	}
}

func (u *mailUsecase) Send(ctx context.Context, req mail.MailSendRequest) {
	span, _ := opentracing.StartSpanFromContext(ctx, "mailUsecase.Send")
	defer span.Finish()

	log.Printf("✉️ Receiving mail.send message: %s", req.To)

	err := u.mailSender.SendEmailWithTemplate(req.To, req.Subject, req.Template, req.Data)
	if err != nil {
		log.Printf("❌ Failed to send email: %v", err)
		return
	}
}
