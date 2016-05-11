package grovepi

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/mattetti/embd"
)

// PinMode indicates if a pin is set as input or output
type PinMode int

const (
	// Out = output
	Out PinMode = iota
	// In = input
	In
)

const (
	//Pins
	A0 byte = 0
	A1 byte = 1
	A2 byte = 2

	D2 byte = 2
	D3 byte = 3
	D4 byte = 4
	D5 byte = 5
	D6 byte = 6
	D7 byte = 7
	D8 byte = 8

	//Cmd format
	DIGITAL_READ  = 1
	DIGITAL_WRITE = 2
	ANALOG_READ   = 3
	ANALOG_WRITE  = 4
	PIN_MODE      = 5
	DHT_READ      = 40

	I2C_SMBUS_READ           = 1
	I2C_SMBUS_WRITE          = 0
	I2C_SMBUS_BYTE_DATA      = 2
	I2C_SMBUS_I2C_BLOCK_DATA = 8
	I2C_SMBUS_BLOCK_MAX      = 32

	// Talk to bus
	I2C_SMBUS = 0x0720

	// Set bus slave
	I2C_SLAVE = 0x0703
)

type Grovepi struct {
	Bus  embd.I2CBus
	Addr byte

	mu sync.RWMutex
}

func New(bus embd.I2CBus) *Grovepi {
	return &Grovepi{
		Bus:  bus,
		Addr: 0x4,
	}
}

// Close stops the controller and resets mode and pwm controller registers.
func (g *Grovepi) Close() error {

	glog.V(1).Infof("Grovepi: close request received")
	return g.Bus.Close()
}

func (g *Grovepi) PinMode(pin byte, mode PinMode) error {
	var b []byte
	var modeB byte = 1
	if mode == In {
		modeB = 0
	}

	b = []byte{PIN_MODE, pin, modeB, 0}
	err := g.Bus.WriteBytes(g.Addr, b)
	if err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (g *Grovepi) DigitalWrite(pin byte, val byte) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	b := []byte{DIGITAL_WRITE, pin, val, 0}
	fmt.Println(g.Addr, b)
	err := g.Bus.WriteBytes(g.Addr, b)
	if err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)

	return nil
}

/*
device.module.Lock()
	defer device.module.Unlock()

	device.sendSlaveAddress()

	buffer := make([]byte, len(data)+1)
	buffer[0] = byte(len(data))
	copy(buffer[1:], data)

	//	buffer := make([]byte, numBytes+2)

	busData := i2c_smbus_ioctl_data{
		read_write: I2C_SMBUS_WRITE,
		command:    command,
		size:       I2C_SMBUS_I2C_BLOCK_DATA,
		data:       uintptr(unsafe.Pointer(&buffer[0])),
	}

	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(device.module.fd.Fd()), I2C_SMBUS, uintptr(unsafe.Pointer(&busData)))
	if err != 0 {
		return syscall.Errno(err)
	}

	return nil
*/
