package liftio

/*
CFLAGS = -std=c99 -g -Wall -O2 -I . -MMD
LDFLAGS = -lpthread -lcomedi -g -lm
#include "io.h"
*/
import "C"
import "log"

/**
  Initialize libComedi in "Sanntidssalen"
*/

func io_init() bool {
	n, err := C.io_init()
	if err != nil {
		log.Fatal("Error interfacing C driver: ", err)
	}
	return bool(n)
}

/**
  Sets a digital channel bit.
*/
func io_set_bit(channel int) {
	_, err := C.io_set_bit(C.int(channel))
	if err != nil {
		log.Fatal("Error interfacing C driver: ", err)
	}
}

/**
  Clears a digital channel bit.
*/
func io_clear_bit(channel int) {
	_, err := C.io_clear_bit(C.int(channel))
	if err != nil {
		log.Fatal("Error interfacing C driver: ", err)
	}
}

/**
  Writes a value to an analog channel.
*/
func io_write_analog(channel, value int) {
	_, err := C.io_write_analog(C.int(channel), C.int(value))
	if err != nil {
		log.Fatal("Error interfacing C driver: ", err)
	}
}

/**
  Reads a bit value from a digital channel.
*/
func io_read_bit(channel int) bool {
	n, err := C.io_read_bit(C.int(channel))
	if err != nil {
		log.Fatal("Error interfacing C driver: ", err)
	}
	return bool(n)
}

/**
  Reads a bit value from an analog channel.
*/
func io_read_analog(channel int) {
	n, err := C.io_read_analog(C.int(channel))
	if err != nil {
		log.Fatal("Error interfacing C driver: ", err)
	}
	return int(n)
}
