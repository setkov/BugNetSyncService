package Common

type MessengerService interface {
	SendMessage(message string) error
}
