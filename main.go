package main

import (
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
	time.Sleep(time.Millisecond * 250)
	roomba.DriveStop()
	roomba.Stop()
	os.Exit(0)
}

var roomba *oibot.OIBot

func main() {

	go quitme()

	roombaInfoLog := log.New(os.Stdout, "", log.LstdFlags)
	roombaErrLog := log.New(os.Stdout, "", log.LstdFlags)
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
	h1move := device.OnMove(1)
	h2move := device.OnMove(3)

	// var hpos Joysticks.CoordsEvent
	var left, right int16

	// start feeding OS events onto the event channels.
	go device.ParcelOutEvents()

	go func() {
		for {
			select {
			case <-b1press:
				log.Println("button #1 pressed")
			case h := <-h1move:
				hpos := h.(Joysticks.CoordsEvent)
				left = int16(hpos.Y * -300)
			case h := <-h2move:
				hpos := h.(Joysticks.CoordsEvent)
				right = int16(hpos.X * -300)
			}
		}
	}()

	var lastleft, lastright int16

	for {

		// log.Println(left, right)

		if left != lastleft || right != lastright {
			roomba.DriveWheels(right, left)
		}

		if lastleft != 0 && lastright != 0 && left == 0 && right == 0 {
			log.Println("stop")
			roomba.DriveStop()
		}

		lastleft = left
		lastright = right

		// if left != 0 && right != 0 && (left != lastleft || right != lastright) {
		// 	if left != lastleft {
		// 		lastleft = left
		// 	}
		// 	if right != lastright {
		// 		lastright = right
		// 	}
		// 	log.Println(lastleft, lastright)
		// 	roomba.Drive(lastleft, lastright)
		// }

		// if speed != lastspeed && speed == 0 {
		// 	lastspeed = 0
		// 	lastradius = 0
		// 	log.Println("stop")
		// 	roomba.DriveStop()
		// }

		time.Sleep(time.Millisecond * 100)
	}

}
