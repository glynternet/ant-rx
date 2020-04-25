package main

import (
	"fmt"

	"github.com/half2me/antgo/message"
)

type AntDeviceMessageHandler interface {
	SpeedAndCadenceMessage(message.SpeedAndCadenceMessage) error
	PowerMessage(message.PowerMessage) error
	HeartRateMessage(HeartRateMessage) error
	Unknown(string, message.AntBroadcastMessage) error
}

func deviceMessageHandler(h AntDeviceMessageHandler) func(m message.AntBroadcastMessage) error {
	decodeClass := deviceClassDecoder(deviceClasses())
	return func(m message.AntBroadcastMessage) error {
		md := m.DeviceType()
		c, err := decodeClass(md)
		if err != nil {
			return err
		}

		switch c {
		case deviceClassBikeSpeedAndCadenceSensor:
			return h.SpeedAndCadenceMessage(message.SpeedAndCadenceMessage(m))
		case deviceClassBikePowerSensor:
			return h.PowerMessage(message.PowerMessage(m))
		case deviceClassHeartRateSensor:
			return h.HeartRateMessage(HeartRateMessage(m))
		}
		return h.Unknown(c, m)
	}
}

func deviceClassDecoder(classes map[byte]string) func(b byte) (string, error) {
	return func(b byte) (string, error) {
		class, ok := classes[b]
		if !ok {
			return deviceClassUnknown, unknownDeviceError(b)
		}
		return class, nil
	}
}

type unknownDeviceError byte

func (dev unknownDeviceError) Error() string {
	return fmt.Sprintf("unknown device: %X", byte(dev))
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
