package services

// unified interface for all service layer implementation notifications
type Notifier[T any] interface {
	Send(req T) error
}
