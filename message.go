package main

import (
	"github.com/half2me/antgo/message"
)

type Visitor interface {
	SpeedAndCadenceMessage(message.SpeedAndCadenceMessage) error
	PowerMessage(message.PowerMessage) error
	Unknown(message.AntBroadcastMessage) error
}

func VisitMessage(v Visitor, m message.AntBroadcastMessage) error {
	switch m.DeviceType() {
	case message.DEVICE_TYPE_SPEED_AND_CADENCE:
		return v.SpeedAndCadenceMessage(message.SpeedAndCadenceMessage(m))
	case message.DEVICE_TYPE_POWER:
		return v.PowerMessage(message.PowerMessage(m))
	}
	return v.Unknown(m)
}
