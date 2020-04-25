package ant

import (
	"github.com/half2me/antgo/message"
	"github.com/pkg/errors"
)

type PacketHandler func(message.AntPacket) error

func NewPacketHandler(handler MessageHandler) PacketHandler {
	decodePacketClass := PacketClassDecoder(PacketClasses())
	return func(packet message.AntPacket) error {
		class, err := decodePacketClass(packet.Class())
		if err != nil {
			return err
		}
		switch class {
		case PacketClassBroadcastData:
			if err := handler.BroadcastMessage(message.AntBroadcastMessage(packet)); err != nil {
				return errors.Wrap(err, "visiting message,")
			}
		default:
			if err := handler.Unknown(class, packet); err != nil {
				return errors.Wrap(err, "handling unknown packet class")
			}
		}
		return nil
	}
}

type MessageHandler interface {
	BroadcastMessage(message.AntBroadcastMessage) error
	Unknown(string, message.AntPacket) error
}

func NewMessageHandler(h AntDeviceMessageHandler) func(m message.AntBroadcastMessage) error {
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

type AntDeviceMessageHandler interface {
	SpeedAndCadenceMessage(message.SpeedAndCadenceMessage) error
	PowerMessage(message.PowerMessage) error
	HeartRateMessage(HeartRateMessage) error
	Unknown(string, message.AntBroadcastMessage) error
}

type OptionalPacketHandlers struct {
	BroadcastMessageHandler func(message.AntBroadcastMessage) error
	UnknownHandler          func(string, message.AntPacket) error
}

func (p OptionalPacketHandlers) BroadcastMessage(msg message.AntBroadcastMessage) error {
	if p.BroadcastMessageHandler == nil {
		return nil
	}
	return p.BroadcastMessageHandler(msg)
}

func (p OptionalPacketHandlers) Unknown(class string, packet message.AntPacket) error {
	if p.UnknownHandler == nil {
		return nil
	}
	return p.UnknownHandler(class, packet)
}
