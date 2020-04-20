package main

import (
	"fmt"

	"github.com/half2me/antgo/message"
)

type printer struct{}

func (p printer) SpeedAndCadenceMessage(message message.SpeedAndCadenceMessage) error {
	fmt.Println(message)
	return nil
}

func (p printer) PowerMessage(message message.PowerMessage) error {
	fmt.Println(message)
	return nil
}

func (p printer) Unknown(message message.AntBroadcastMessage) error {
	fmt.Println(message)
	return nil
}
