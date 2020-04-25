package main

import (
	"encoding/json"
	"fmt"

	"github.com/glynternet/ant-rx/ant"
	"github.com/half2me/antgo/message"
	"github.com/pkg/errors"
)

func newPacketPrinter(printUnknown bool) ant.MessageHandler {
	return ant.OptionalPacketHandlers{
		UnknownHandler: func(class string, p message.AntPacket) error {
			return jsonPrint("unknown packet class", struct {
				Class string
			}{
				Class: class,
			})
		},
		BroadcastMessageHandler: ant.NewMessageHandler(deviceMessagePrinter{printUnknown: printUnknown}),
	}
}

type deviceMessagePrinter struct {
	printUnknown    bool
	prevSpeedCadMsg message.SpeedAndCadenceMessage
}

func (p deviceMessagePrinter) SpeedAndCadenceMessage(msg message.SpeedAndCadenceMessage) error {
	type speed struct {
		EventTime                 uint16
		CumulativeRevolutionCount uint16
	}
	type cadence struct {
		EventTime                 uint16
		CumulativeRevolutionCount uint16
	}
	type speedAndCadence struct {
		Speed   speed
		Cadence cadence
	}

	return jsonPrintDeviceMessage(message.AntBroadcastMessage(msg).DeviceNumber(), "speed and cadence", speedAndCadence{
		Speed: speed{
			EventTime:                 msg.SpeedEventTime(),
			CumulativeRevolutionCount: msg.CumulativeSpeedRevolutionCount(),
		},
		Cadence: cadence{
			EventTime:                 msg.CadenceEventTime(),
			CumulativeRevolutionCount: msg.CumulativeCadenceRevolutionCount(),
		},
	})
}

func (p deviceMessagePrinter) PowerMessage(msg message.PowerMessage) error {
	type power struct {
		AccumulatedPower     uint16
		EventCount           uint8
		InstantaneousCadence uint8
		InstantaneousPower   uint16
	}

	return jsonPrintDeviceMessage(message.AntBroadcastMessage(msg).DeviceNumber(), "power", power{
		AccumulatedPower:     msg.AccumulatedPower(),
		EventCount:           msg.EventCount(),
		InstantaneousCadence: msg.InstantaneousCadence(),
		InstantaneousPower:   msg.InstantaneousPower(),
	})
}

func (p deviceMessagePrinter) HeartRateMessage(msg ant.HeartRateMessage) error {
	msg.HeartRate()
	msg.BeatCount()
	type heartRate struct {
		Rate  uint8
		Count uint8
	}
	return jsonPrintDeviceMessage(message.AntBroadcastMessage(msg).DeviceNumber(), "heart rate", heartRate{
		Rate:  msg.HeartRate(),
		Count: msg.BeatCount(),
	})
}

func (p deviceMessagePrinter) Unknown(s string, message message.AntBroadcastMessage) error {
	return jsonPrintDeviceMessage(message.DeviceNumber(), "unknown", struct {
		Class string
	}{
		Class: s,
	})
}

func jsonPrintDeviceMessage(deviceNumber uint16, name string, v interface{}) error {
	return jsonPrint(name, devicePrintMessage{
		DeviceNumber: deviceNumber,
		Message:      v,
	})
}

func jsonPrint(name string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return errors.Wrapf(err, "json marshalling %s", name)
	}
	fmt.Println(string(data))
	return nil
}

type devicePrintMessage struct {
	DeviceNumber uint16
	Message      interface{}
}
