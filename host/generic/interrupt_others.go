//+build darwin,!linux

// mock implementation so the package can properly compile on non
// linux machines during dev.
package generic

import "github.com/mattetti/embd"

func initEpollListener() *epollListener {
	return nil
}

func registerInterrupt(pin *digitalPin, handler func(embd.DigitalPin)) error {
	return nil
}

func unregisterInterrupt(pin *digitalPin) error {
	return nil
}
