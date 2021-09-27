package models

// Controller is interface of Archer "brain"
type Controller interface {
	MessageTo(neighbor Archer)
	Fire()
	Start()
	Stop() error
}
