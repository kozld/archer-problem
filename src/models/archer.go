package models

import (
	"fmt"
)

///////////////////////////////////////
//////////// ARCHER MODEL  ////////////
///////////////////////////////////////

// NewArcher can construct new instance of Archer
func NewArcher(controller Controller) *Archer {
	return &Archer{make(map[string]interface{}), controller}
}

// Archer represents Archer
type Archer struct {
	memory map[string]interface{}
	cpu    Controller
	//Host   host.Host
	//cancel context.CancelFunc
}

// MessageTo can send message to neighbor Archer
func (a *Archer) MessageTo(neighbor Archer) {
	fmt.Println("I'm sending message")
	a.cpu.MessageTo(neighbor)
}

// Fire can to fire
func (a *Archer) Fire() {
	fmt.Println("I'm fire!")
	a.cpu.Fire()
}

// DoNothing can do nothing
func (a *Archer) DoNothing() {
	fmt.Println("Do nothing...")
}

// GetFromMemory can remember data from Archer memory
func (a *Archer) GetFromMemory(key string) interface{} {
	fmt.Printf("I remember data from memory")
	return a.memory[key]
}

// SaveToMemory can remember data in Archer memory
func (a *Archer) SaveToMemory(key string, value interface{}) {
	fmt.Printf("I remember data in memory")
	a.memory[key] = value
}

// Burn can burn Archer
func (a *Archer) Burn() {
	a.cpu.Start()
}

// Destroy can destroy Archer
func (a *Archer) Destroy() {
	fmt.Println("I self destruct :(")
	a.cpu.Stop()
	//a.Host.Close()
}
