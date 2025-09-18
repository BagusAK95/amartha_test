package messaging

import (
	"context"

	"github.com/BagusAK95/amarta_test/internal/domain/mail"
)

type mailHandler struct {
	usecase mail.IMailUsecase
}

func NewMailHandler(usecase mail.IMailUsecase) *mailHandler {
	return &mailHandler{
		usecase: usecase,
	}
}

func (h *mailHandler) Send(msg mail.MailSendRequest) {
	ctx := context.Background()

	h.usecase.Send(ctx, msg)
}
