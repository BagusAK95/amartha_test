package usecase

import (
	"context"
	"log"

	"github.com/BagusAK95/amarta_test/internal/domain/mail"
	mailsender "github.com/BagusAK95/amarta_test/internal/infrastructure/mail"
	"go.opentelemetry.io/otel"
)

var tracerName = "MailUsecase"
var tracer = otel.Tracer(tracerName)

type mailUsecase struct {
	mailSender mailsender.ISender
}

func NewMailUsecase(mailSender mailsender.ISender) mail.IMailUsecase {
	return &mailUsecase{
		mailSender: mailSender,
	}
}

func (u *mailUsecase) Send(ctx context.Context, req mail.MailSendRequest) {
	_, span := tracer.Start(ctx, tracerName+".Send")
	defer span.End()

	log.Printf("✉️ Receiving mail.send message: %s", req.To)

	err := u.mailSender.SendEmailWithTemplate(req.To, req.Subject, req.Template, req.Data)
	if err != nil {
		log.Printf("❌ Failed to send email: %v", err)
		return
	}
}
