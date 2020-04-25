package main

import (
	"fmt"

	"github.com/half2me/antgo/message"
)

type packetHandler struct {
	broadcastMessage func(message.AntBroadcastMessage) error
	unhandled        func(string, message.AntPacket) error
}

func (p packetHandler) BroadcastMessage(msg message.AntBroadcastMessage) error {
	if p.broadcastMessage == nil {
		return nil
	}
	return p.broadcastMessage(msg)
}

func (p packetHandler) Unknown(class string, packet message.AntPacket) error {
	if p.unhandled == nil {
		return nil
	}
	return p.unhandled(class, packet)
}
func newPacketPrinter(printUnknown bool, packetClasses, deviceClasses map[byte]string) packetHandler {
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

func (p deviceMessagePrinter) Unknown(message message.AntBroadcastMessage) error {
	dt, known := deviceClasses()[message.DeviceType()]
	if known {
		fmt.Printf("%s %s\n", dt, message)
		return nil
	}
	if !p.printUnknown {
		return nil
	}
	fmt.Printf("Unknown Device Type: %d, %s\n", message.DeviceType(), message)
	return nil
}
