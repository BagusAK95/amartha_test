package mail

type MailSendRequest struct {
	To       string
	Subject  string
	Template string
	Data     map[string]any
}
