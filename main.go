package main

import (
	"fmt"
    "math"
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
)


type Color struct {
    pin rpio.Pin
    intensity chan int
    quit chan int
}

func NewColor(pin rpio.Pin) *Color {
    c := new(Color)
    c.pin = pin
    c.pin.Output()
    //go c.Pwm(c.intensity, c.quit, 1000, 0)
    return c
}

func (c *Color) Pwm(intensity, quit chan int, freq int, duty float64) {
    samples := freq * 100
    for {
        select {
            case val := <-intensity:
                fmt.Println("received", val)
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

    red := NewColor(22)

    for x := 0; x < 4; x++ {
        red.pin.Toggle()
        time.Sleep(time.Second / 10)
    }
    //for x := 0; x < 255; x++ {
        //red.intensity <- x
        //time.Sleep(time.Second * 5 / 255)
    //}
}
