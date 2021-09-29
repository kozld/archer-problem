package models

///////////////////////////////////////
// CONTROLLER INTERFACE
///////////////////////////////////////

// Controller is interface of Archer brain
type Controller interface {
	Message(to Archer, message interface{})
	Fire()
	Start()
	Stop() error
}
