package ant

import (
	"fmt"

	"github.com/half2me/antgo/message"
)

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
