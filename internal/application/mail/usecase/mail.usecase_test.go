package usecase_test

import (
	"context"
	"testing"

	"github.com/BagusAK95/amarta_test/internal/application/mail/usecase"
	"github.com/BagusAK95/amarta_test/internal/domain/mail"
	mailMock "github.com/BagusAK95/amarta_test/internal/infrastructure/mail/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSend(t *testing.T) {
	ctx := context.Background()
	req := mail.MailSendRequest{
		To:      "test@example.com",
		Subject: "test subject",
	}

	t.Run("success", func(t *testing.T) {
		mailSender := new(mailMock.MockISender)
		mailSender.On("SendEmailWithTemplate", req.To, req.Subject, mock.AnythingOfType("string"), mock.Anything).Return(nil)

		uc := usecase.NewMailUsecase(mailSender)
		uc.Send(ctx, req)

		mailSender.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mailSender := new(mailMock.MockISender)
		mailSender.On("SendEmailWithTemplate", req.To, req.Subject, mock.AnythingOfType("string"), mock.Anything).Return(assert.AnError)

		uc := usecase.NewMailUsecase(mailSender)
		uc.Send(ctx, req)

		mailSender.AssertExpectations(t)
	})
}
