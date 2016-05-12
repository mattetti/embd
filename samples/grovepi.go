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
	gp *grovepi.Grovepi
)

// Connect the LED to D2
// This LED will blink for 5 seconds and the code will exit
func main() {
	flag.Parse()

	bus := embd.NewI2CBus(1)
	gp = grovepi.New(bus)
	defer gp.Close()

	ledpin := grovepi.D2
	err := gp.PinMode(ledpin, grovepi.Out)
	if err != nil {
		panic(err)
	}
	// cleanup
	defer func() {
		gp.DigitalWrite(ledpin, 0)
	}()

	stopC := make(chan bool)
	go func() {
		blink(ledpin, stopC)
	}()

	time.Sleep(5 * time.Second)
	stopC <- true
}

func blink(ledpin byte, stop chan bool) {
	var state uint8 = 1
	for {
		select {
		case <-stop:
			return
		default:
			if err := gp.DigitalWrite(ledpin, byte(state%2)); err != nil {
				fmt.Println(err)
			}
			state++
			time.Sleep(500 * time.Millisecond)
		}
	}
}
