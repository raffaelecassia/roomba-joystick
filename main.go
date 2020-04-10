package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	Joysticks "github.com/splace/joysticks"

	"github.com/ardnew/oibot"
)

func quitme() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("QUIT")
	roomba.DriveStop()
	roomba.Stop()
	os.Exit(0)
}

var roomba *oibot.OIBot

const opcPWMMotors oibot.OpCode = 144

func main() {

	go quitme()

	// roombaInfoLog := log.New(os.Stdout, "", log.LstdFlags)
	roombaInfoLog := log.New(ioutil.Discard, "", log.LstdFlags)
	roombaErrLog := log.New(os.Stdout, "ROOMBA ERROR! ", log.LstdFlags)

	roomba = oibot.MakeOIBot(roombaInfoLog, roombaErrLog, false, "/dev/ttyUSB0", 115200, oibot.DefaultReadTimeoutMS)
	time.Sleep(time.Second * 1)
	roomba.Safe()

	// try connecting to specific controller.
	// the index is system assigned, typically it increments on each new controller added.
	// indexes remain fixed for a given controller, if/when other controller(s) are removed.
	device := Joysticks.Connect(1)
	if device == nil {
		panic("no HIDs")
	}

	// using Connect allows a device to be interrogated
	log.Printf("HID#1:- Buttons:%d, Hats:%d\n", len(device.Buttons), len(device.HatAxes)/2)

	// get/assign channels for specific events
	b1press := device.OnClose(1)
	b3press := device.OnClose(3)
	bOptionpress := device.OnClose(10)
	h1move := device.OnMove(1)
	h2move := device.OnMove(2)

	var hpos Joysticks.CoordsEvent
	var hpos2 Joysticks.CoordsEvent

	// start feeding OS events onto the event channels.
	go device.ParcelOutEvents()

	brushes := false

	go func() {
		for {
			select {
			case <-b1press:
				log.Println("button #1 pressed")
				if brushes {
					roomba.Write(opcPWMMotors, uint8(0), uint8(0), uint8(0))
					brushes = false
				} else {

					roomba.Write(opcPWMMotors, uint8(127), uint8(127), uint8(127))
					brushes = true
				}

			case <-b3press:
				log.Println("button #3 pressed")
				roomba.SeekDock()

			case <-bOptionpress:
				log.Println("RESET")
				roomba.DriveStop()
				time.Sleep(time.Millisecond * 100)
				roomba.Stop()
				time.Sleep(time.Millisecond * 500)
				roomba.Start()
				time.Sleep(time.Millisecond * 500)
				roomba.Safe()

			case h := <-h1move:
				hpos = h.(Joysticks.CoordsEvent)

			case h := <-h2move:
				hpos2 = h.(Joysticks.CoordsEvent)
			}
		}
	}()

	var speed, lastspeed, radius, lastradius, tmp, sign int16
	var lastrotate float32

	for {
		tmp = int16(hpos2.Y * 100) // [0, +100]
		if tmp != 0 {
			if tmp > 0 {
				sign = -1
			} else {
				sign = 1
				tmp = tmp * -1
			}
			radius = int16(mappete(int(tmp), 0, 100, 250, 100)) * sign
		} else {
			radius = 0
		}

		speed = int16(hpos.Y * -500)

		if speed != 0 && (speed != lastspeed || radius != lastradius) {
			if speed != lastspeed {
				lastspeed = speed
			}
			if radius != lastradius {
				lastradius = radius
			}
			// log.Println(lastspeed, lastradius)
			roomba.Drive(lastspeed, lastradius)
		}

		if speed != lastspeed && speed == 0 {
			lastspeed = 0
			lastradius = 0
			log.Println("stop")
			roomba.DriveStop()
		}

		if speed == 0 && lastrotate != hpos2.Y {
			lastrotate = hpos2.Y
			l := int16(hpos2.Y * 500)
			r := int16(hpos2.Y * -500)
			log.Println("ROTATE", l, r)
			roomba.DriveWheels(r, l)

		}

		time.Sleep(time.Millisecond * 150)
	}

}

func mappete(x int, inmin int, inmax int, outmin int, outmax int) int {
	return (x-inmin)*(outmax-outmin)/(inmax-inmin) + outmin
}
