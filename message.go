package main

import (
	"fmt"

	"github.com/half2me/antgo/message"
)

type AntBroadcastMessageVisitor interface {
	SpeedAndCadenceMessage(message.SpeedAndCadenceMessage) error
	PowerMessage(message.PowerMessage) error
	HeartRateMessage(HeartRateMessage) error
	Unknown(message.AntBroadcastMessage) error
}

type unknownDeviceError byte

func (dev unknownDeviceError) Error() string {
	return fmt.Sprintf("unknown device: %X", byte(dev))
}

func VisitMessage(v AntBroadcastMessageVisitor, m message.AntBroadcastMessage) error {
	md := m.DeviceType()
	// TODO(glynternet): non-global way of not calling deviceClasses every invocation
	c, ok := deviceClasses()[md]
	if !ok {
		return unknownDeviceError(md)
	}

	switch c {
	case deviceClassBikeSpeedAndCadenceSensor:
		return v.SpeedAndCadenceMessage(message.SpeedAndCadenceMessage(m))
	case deviceClassBikePowerSensor:
		return v.PowerMessage(message.PowerMessage(m))
	case deviceClassHeartRateSensor:
		return v.HeartRateMessage(HeartRateMessage(m))
	}
	return v.Unknown(m)
}

type HeartRateMessage message.AntPacket

func (hrm HeartRateMessage) BeatCount() uint8 {
	return message.AntBroadcastMessage(hrm).Content()[6]
}

func (hrm HeartRateMessage) HeartRate() uint8 {
	return message.AntBroadcastMessage(hrm).Content()[7]
}

func (hrm HeartRateMessage) String() string {
	return fmt.Sprintf("HRM: Rate:%d, Count:%d", hrm.HeartRate(), hrm.BeatCount())
}
