# roomba-joystick

A quarantine mini project about an iRobot Roomba 620 and a PS4 controller.

[video example](https://www.youtube.com/watch?v=GZHvjnbsxjY)


## Ingredients
* Roomba 6xx
* PS4 Controller
* Something that can run programs (with WiFi and Bluetooth, like a RaspberryPi3)
* A way to connect the 0-5V levels serial port to the raspberry

Any roomba model that has the mini-din connector in his back and is supported by the iRobot Create2 Open Interface Specs can be used.

PS4 controller is not mandatory. Any HID controller connected either via BT or USB will work (experience with axis and buttons may vary, code edits may be required).

To power things, Roomba mini-din has an unregulated battery pin, but it requires additional hardware. So I used a powerbank instead.

There are at least 2 ways to have a serial port on a raspberry. A usb-to-serial adapter that works with 0-5V levels or use GPIO 14 (TXD) and 15 (RXD) with a level shifter (raspberry only works with 3.3V). 
Connect TXD-RXD two times and you are ready to go.


## FIRST THING FIRST
Always double-check voltages, polarity, connections, documentation, demoniac presence, etc, before powering things up. You may burn your roomba, your raspberry, or your home.


## How to run all of this

* compile: `GOOS=linux GOARCH=arm GOARM=7 go build`
* copy `roomba-joystick` to raspberry
* ssh raspberry
* `sudo bluetoothctl` to open a sort of bt shell
* put your bt controller in pairing mode
* `scan on`
* wait for something like "`[NEW] Device $BLUETOOTH_DEVICE_ADDRESS Wireless Controller`"
* type these:
```
pair $BLUETOOTH_DEVICE_ADDRESS
trust $BLUETOOTH_DEVICE_ADDRESS
connect $BLUETOOTH_DEVICE_ADDRESS
exit
```
* execute `./roomba-joystick`
* ???
* profit
