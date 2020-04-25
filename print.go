package main

import (
	"fmt"

	"github.com/half2me/antgo/message"
)

func newPacketPrinter(printUnknown bool) packetHandler {
	return packetHandler{
		unknown: func(class string, p message.AntPacket) error {
			return printf("Received unknown packet: %s\n", class)
		},
		broadcastMessage: deviceMessageHandler(deviceMessagePrinter{printUnknown: printUnknown}),
	}
}

type deviceMessagePrinter struct {
	printUnknown bool
}

func (p deviceMessagePrinter) SpeedAndCadenceMessage(message message.SpeedAndCadenceMessage) error {
	return print(message)
}

func (p deviceMessagePrinter) PowerMessage(message message.PowerMessage) error {
	return print(message)
}

func (p deviceMessagePrinter) HeartRateMessage(message HeartRateMessage) error {
	return print(message)
}

func (p deviceMessagePrinter) Unknown(s string, message message.AntBroadcastMessage) error {
	return printf("Unknown Device Type: %d, %s\n", s, message)
}

func print(v interface{}) error {
	fmt.Println(v)
	return nil
}

func printf(format string, a ...interface{}) error {
	fmt.Printf(format, a...)
	return nil
}
