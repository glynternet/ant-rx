package main

import (
	"fmt"

	"github.com/half2me/antgo/message"
)

func newPacketPrinter(printUnknown bool) packetHandler {
	return packetHandler{
		unhandled: func(class string, p message.AntPacket) error {
			fmt.Printf("Received unhandled packet: %s\n", class)
			return nil
		},
		broadcastMessage: deviceMessageHandler(deviceMessagePrinter{printUnknown: printUnknown}),
	}
}

type deviceMessagePrinter struct {
	printUnknown bool
}

func (p deviceMessagePrinter) SpeedAndCadenceMessage(message message.SpeedAndCadenceMessage) error {
	fmt.Println(message)
	return nil
}

func (p deviceMessagePrinter) PowerMessage(message message.PowerMessage) error {
	fmt.Println(message)
	return nil
}

func (p deviceMessagePrinter) HeartRateMessage(message HeartRateMessage) error {
	fmt.Println(message)
	return nil
}

func (p deviceMessagePrinter) Unknown(s string, message message.AntBroadcastMessage) error {
	fmt.Printf("Unknown Device Type: %d, %s\n", s, message)
	return nil
}
