package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/mattetti/embd"
	"github.com/mattetti/embd/controller/grovepi"
	_ "github.com/mattetti/embd/host/rpi"
)

func main() {
	flag.Parse()

	bus := embd.NewI2CBus(1)
	gp := grovepi.New(bus)
	defer gp.Close()

	pin := grovepi.D2
	err := gp.PinMode(pin, grovepi.Out)
	if err != nil {
		fmt.Println(err)
	}

	for {
		fmt.Println("bleep")
		if err := gp.DigitalWrite(pin, 1); err != nil {
			fmt.Println(err)
		}
		time.Sleep(500 * time.Millisecond)
		if err := gp.DigitalWrite(pin, 0); err != nil {
			fmt.Println(err)
		}
		time.Sleep(500 * time.Millisecond)
	}

}

/*
var g grovepi.GrovePi
	g = *grovepi.InitGrovePi(0x04)
	err := g.PinMode(grovepi.D2, "output")
	if err != nil {
		fmt.Println(err)
	}
	for {
		g.DigitalWrite(grovepi.D2, 1)
		time.Sleep(500 * time.Millisecond)
		g.DigitalWrite(grovepi.D2, 0)
		time.Sleep(500 * time.Millisecond)
	}
*/
