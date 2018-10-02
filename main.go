package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
)


type Color struct {
    pin rpio.Pin
}

func NewColor(pin rpio.Pin) *Color {
    c := new(Color)
    c.pin = pin
    c.pin.Output()
    return c
}

func (c *Color) Pwm(freq int, duty int) {
    t := 1
    cycles := t * freq
    samples := freq * 100
    for i := 0; i < cycles; i++ {
        c.pin.High()
        time.Sleep((time.Second * time.Duration(duty)) / time.Duration(samples))
        c.pin.Low()
        time.Sleep((time.Second * (100 - time.Duration(duty))) / time.Duration(samples))
    }
}

func main() {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()

    red := NewColor(22)

    //for x := 0; x < 4; x++ {
        //red.pin.Toggle()
        //time.Sleep(time.Second / 10)
    //}
    for x := 0; x < 10; x++ {
        red.Pwm(1000, x*10)
    }
}
