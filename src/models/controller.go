package models

// Controller is interface of Archer "brain"
type Controller interface {
	MessageTo(neighbor Archer, message string)
	Fire()
	Start()
	Stop() error
}
