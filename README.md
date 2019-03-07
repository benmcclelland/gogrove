# gogrove

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/benmcclelland/gogrove)

Go library for interacting with GrovePi

Currently only tested with GrovePi firmware version 1.3.0

The Rasberry Pi communicates with GrovePi over I2C.  The following kernel
modules are needed to support this:

* i2c_dev
* i2c_bcm2835

To see if the Rasberry Pi is communicating with the Grove pi, run the following:

```sh
# sudo i2cdetect -y 1
     0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
00:          03 04 -- -- -- -- -- -- -- -- -- -- --
10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
20: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- 3e --
40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
60: -- -- 62 -- -- -- -- -- -- -- -- -- -- -- -- --
70: -- -- -- -- -- -- -- --
```

You should see the "04" for the GrovePi and the "3e" and "62" for the LCD.

If these are not showing up, try reloading the i2c_bcm2835 module:

```sh
# sudo modprobe i2c_dev
# sudo rmmod i2c_bcm2835
# sudo modprobe i2c_bcm2835
```

This package is goroutine safe within a session

Some logic within is based on the Python library aviable [here](https://github.com/DexterInd/GrovePi/tree/master/Software/Python)

Useful links:

* [port description](https://www.dexterindustries.com/GrovePi/engineering/port-description/)
* [some protocol dicsussion](https://www.dexterindustries.com/GrovePi/programming/grovepi-protocol-adding-custom-sensors/)
* [main GrovePi page](https://www.dexterindustries.com/grovepi/)