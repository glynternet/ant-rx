package main

import "github.com/half2me/antgo/message"

type packetHandler struct {
	broadcastMessage func(message.AntBroadcastMessage) error
	unknown          func(string, message.AntPacket) error
}

func (p packetHandler) BroadcastMessage(msg message.AntBroadcastMessage) error {
	if p.broadcastMessage == nil {
		return nil
	}
	return p.broadcastMessage(msg)
}

func (p packetHandler) Unknown(class string, packet message.AntPacket) error {
	if p.unknown == nil {
		return nil
	}
	return p.unknown(class, packet)
}
