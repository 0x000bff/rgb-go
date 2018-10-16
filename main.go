package main

import (
	"fmt"
    "math"
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
)


var (
    redPin = rpio.Pin(22)
    greenPin = rpio.Pin(17)
    bluePin = rpio.Pin(24)
    whitePin = rpio.Pin(27)
)


type ColorComponent struct {
    pin rpio.Pin
    intensity chan int
    quit chan int
}


func NewColorComponent(pin rpio.Pin) *ColorComponent {
    c := new(ColorComponent)
    c.pin = pin
    c.pin.Output()
    c.intensity = make(chan int)
    c.quit = make(chan int)
    go c.Pwm(c.intensity, c.quit, 1000, 0)
    return c
}


func (c *ColorComponent) Pwm(intensity, quit chan int, freq int, duty float64) {
    samples := freq * 100
    for {
        select {
            case val := <-intensity:
                duty = intensityToDutyCycle(val)
                continue
            case <-quit:
                c.pin.Low()
                return
            default:
                c.pin.High()
                time.Sleep((time.Second * time.Duration(duty)) / time.Duration(samples))
                c.pin.Low()
                time.Sleep((time.Second * (100 - time.Duration(duty))) / time.Duration(samples))
        }
    }
}


type Color struct {
    red *ColorComponent
    green *ColorComponent
    blue *ColorComponent
    white *ColorComponent
}


func NewColor() *Color {
    c := new(Color)
    c.red = NewColorComponent(redPin)
    c.green = NewColorComponent(greenPin)
    c.blue = NewColorComponent(bluePin)
    c.white = NewColorComponent(whitePin)
    return c
}


func (c *Color) quit() {
    c.red.quit <- 0
    c.blue.quit <- 0
    c.green.quit <- 0
    c.white.quit <- 0
}


func intensityToDutyCycle(intensity int)  float64 {
    if intensity < 0 {
        intensity = 0
    }
    duty := math.Pow((float64(intensity)/25.5), 2)
    if duty > 100 {
        duty = 100
    }
    return duty
}


func main() {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()

    color := NewColor()

    for x := 1; x < 256; x++ {
        color.red.intensity <- x
        time.Sleep(time.Second / 200)
    }
    color.quit()
}
