//+ build linux
// Generic Interrupt Pins.

package generic

import (
	"fmt"
	"syscall"

	"github.com/mattetti/embd"
)

func initEpollListener() *epollListener {
	fd, err := syscall.EpollCreate1(0)
	if err != nil {
		panic(fmt.Sprintf("Unable to create epoll: %v", err))
	}
	listener := &epollListener{fd: fd, interruptablePins: make(map[int]*interrupt)}

	go func() {
		var epollEvents [MaxGPIOInterrupt]syscall.EpollEvent

		for {
			n, err := syscall.EpollWait(listener.fd, epollEvents[:], -1)
			if err != nil {
				panic(fmt.Sprintf("EpollWait error: %v", err))
			}
			for i := 0; i < n; i++ {
				if irq, ok := listener.interruptablePins[int(epollEvents[i].Fd)]; ok {
					irq.Signal()
				}
			}
		}
	}()
	return listener
}

func registerInterrupt(pin *digitalPin, handler func(embd.DigitalPin)) error {
	l := getEpollListenerInstance()

	pinFd := int(pin.val.Fd())

	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.interruptablePins[pinFd]; ok {
		return ErrorPinAlreadyRegistered
	}

	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN | (syscall.EPOLLET & 0xffffffff) | syscall.EPOLLPRI

	if err := syscall.SetNonblock(pinFd, true); err != nil {
		return err
	}

	event.Fd = int32(pinFd)

	if err := syscall.EpollCtl(l.fd, syscall.EPOLL_CTL_ADD, pinFd, &event); err != nil {
		return err
	}

	l.interruptablePins[pinFd] = &interrupt{pin: pin, handler: handler}

	return nil
}

func unregisterInterrupt(pin *digitalPin) error {
	l := getEpollListenerInstance()

	pinFd := int(pin.val.Fd())

	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.interruptablePins[pinFd]; !ok {
		return nil
	}

	if err := syscall.EpollCtl(l.fd, syscall.EPOLL_CTL_DEL, pinFd, nil); err != nil {
		return err
	}

	if err := syscall.SetNonblock(pinFd, false); err != nil {
		return err
	}

	delete(l.interruptablePins, pinFd)
	return nil
}
