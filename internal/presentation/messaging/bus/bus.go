package bus

import (
	mailhandler "github.com/BagusAK95/amarta_test/internal/application/mail/delivery/messaging"
	"github.com/BagusAK95/amarta_test/internal/domain/mail"
	"github.com/BagusAK95/amarta_test/internal/infrastructure/bus"
)

func NewBusListener(mailBus bus.Bus[mail.MailSendRequest], mailUsecase mail.IMailUsecase) {
	handler := mailhandler.NewMailHandler(mailUsecase)

	mailBus.SubscribeAsync("mail.send", handler.Send, false)
}
