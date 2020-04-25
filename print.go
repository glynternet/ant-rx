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
	printUnknown bool
}

func (p deviceMessagePrinter) SpeedAndCadenceMessage(d ant.Device, msg message.SpeedAndCadenceMessage) error {
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

	return jsonPrintDeviceMessage(d, "speed and cadence", speedAndCadence{
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

func (p deviceMessagePrinter) PowerMessage(d ant.Device, msg message.PowerMessage) error {
	type power struct {
		AccumulatedPower     uint16
		EventCount           uint8
		InstantaneousCadence uint8
		InstantaneousPower   uint16
	}

	return jsonPrintDeviceMessage(d, "power", power{
		AccumulatedPower:     msg.AccumulatedPower(),
		EventCount:           msg.EventCount(),
		InstantaneousCadence: msg.InstantaneousCadence(),
		InstantaneousPower:   msg.InstantaneousPower(),
	})
}

func (p deviceMessagePrinter) HeartRateMessage(d ant.Device, msg ant.HeartRateMessage) error {
	type heartRate struct {
		Rate  uint8
		Count uint8
	}
	return jsonPrintDeviceMessage(d, "heart rate", heartRate{
		Rate:  msg.HeartRate(),
		Count: msg.BeatCount(),
	})
}

func (p deviceMessagePrinter) Unknown(d ant.Device, msg message.AntBroadcastMessage) error {
	if !p.printUnknown {
		return nil
	}
	return jsonPrintDeviceMessage(d, "unknown", struct {
		Message string
	}{
		Message: string(msg),
	})
}

func jsonPrintDeviceMessage(device ant.Device, name string, v interface{}) error {
	return jsonPrint(name, printableDeviceMessage{
		Device:  device,
		Message: v,
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

type printableDeviceMessage struct {
	Device  ant.Device
	Message interface{}
}
