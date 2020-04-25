package main

import (
	"fmt"

	"github.com/half2me/antgo/message"
)

type printer struct {
	printUnknown bool
}

func (p printer) SpeedAndCadenceMessage(message message.SpeedAndCadenceMessage) error {
	fmt.Println(message)
	return nil
}

func (p printer) PowerMessage(message message.PowerMessage) error {
	fmt.Println(message)
	return nil
}

func (p printer) HeartRateMessage(message HeartRateMessage) error {
	fmt.Println(message)
	return nil
}

func (p printer) Unknown(message message.AntBroadcastMessage) error {
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
