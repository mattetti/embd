package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/mattetti/embd"
	"github.com/mattetti/embd/controller/grovepi"
	_ "github.com/mattetti/embd/host/rpi"
)

var (
	gp     *grovepi.Grovepi
	ledpin = grovepi.D2
	btnpin = grovepi.D3

	isBlinking bool
	blinkDelay = 80 * time.Millisecond
)

// Connect a LED to D2
// Connect a button to D3
// Press the button to toggle the LED between blinking and not
func main() {
	flag.Parse()

	bus := embd.NewI2CBus(1)
	gp = grovepi.New(bus)
	defer gp.Close()

	// LED
	err := gp.PinMode(ledpin, grovepi.Out)
	if err != nil {
		panic(err)
	}

	// Button
	err = gp.PinMode(btnpin, grovepi.In)
	if err != nil {
		panic(err)
	}

	// channel used to communicate with the blinking
	// code
	stopC := make(chan bool)

	// cleanup
	defer func() {
		if isBlinking {
			stopC <- true
		}
	}()

	// TODO: try to use interupts instead
	for {
		v, err := gp.DigitalRead(btnpin)
		if err != nil {
			fmt.Println(err)
		}

		if v == 1 {
			fmt.Println("button pressed")
			if isBlinking {
				stopC <- true
			} else {
				go func() {
					blink(ledpin, stopC)
				}()
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}

	}

}

func blink(ledpin byte, stop chan bool) {
	var state uint8 = 1
	isBlinking = true

	defer func() {
		isBlinking = false
		gp.DigitalWrite(ledpin, 0)
	}()

	for {
		select {
		case <-stop:
			return
		default:
			if err := gp.DigitalWrite(ledpin, byte(state%2)); err != nil {
				fmt.Println(err)
			}
			state++
			fmt.Println(state)
			time.Sleep(blinkDelay)
		}
	}
}
