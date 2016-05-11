// Generic Interrupt Pins.

package generic

import (
	"errors"
	"sync"

	"github.com/mattetti/embd"
)

const (
	MaxGPIOInterrupt = 64
)

var ErrorPinAlreadyRegistered = errors.New("pin interrupt already registered")

type interrupt struct {
	pin            embd.DigitalPin
	initialTrigger bool
	handler        func(embd.DigitalPin)
}

func (i *interrupt) Signal() {
	if !i.initialTrigger {
		i.initialTrigger = true
		return
	}
	i.handler(i.pin)
}

type epollListener struct {
	mu                sync.Mutex // Guards the following.
	fd                int
	interruptablePins map[int]*interrupt
}

var epollListenerInstance *epollListener

func getEpollListenerInstance() *epollListener {
	if epollListenerInstance == nil {
		epollListenerInstance = initEpollListener()
	}
	return epollListenerInstance
}
