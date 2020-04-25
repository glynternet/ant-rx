package ant

import "fmt"

const (
	packetClassUnassignChannel  = "unassign_channel"
	packetClassAssignChannel    = "assign_channel"
	packetClassChannelID        = "channel_id"
	packetClassChannelPeriod    = "channel_period"
	packetClassSearchTKL        = "search_timeout"
	packetClassChannelFrequency = "channel_frequency"
	packetClassSetNetwork       = "set_network"
	packetClassTxPower          = "tx_power"
	packetClassIdListAdd        = "id_list_add"
	packetClassIdListConfig     = "id_list_config"
	packetClassChanneltxPower   = "channel_tx_power"
	packetClassLpSearchTimeout  = "lp_search_timeout"
	packetClassSetSerialNumber  = "set_serial_number"
	packetClassEnableExtMsgs    = "enable_ext_msgs"
	packetClassEnableLED        = "enable_led"
	packetClassSystemReset      = "system_reset"
	packetClassOpenChannel      = "open_channel"
	packetClassCloseChannel     = "close_channel"
	packetClassOpenRXScanCH     = "open_rx_scan_ch"
	packetClassReqMessage       = "req_message"
	PacketClassBroadcastData    = "broadcast_data"
	packetClassAckData          = "ack_data"
	packetClassBurstData        = "burst_data"
	packetClassChannelEvent     = "channel_event"
	packetClassChannelStatus    = "channel_status"
	packetClassVersion          = "version"
	packetClassCapabilities     = "capabilities"
	packetClassSerialNumber     = "serial_number"
	packetClassNotifStartup     = "notif_startup"
	packetClassCwInit           = "cw_init"
	packetClassCwTest           = "cw_test"
	packetClassUnknown          = "UnknownHandler"
)

func PacketClasses() map[byte]string {
	return map[byte]string{
		// From here: https://github.com/GoldenCheetah/GoldenCheetah/blob/3a31f5d131df46c90e25810a876ee4c5e0db5512/src/ANT/ANT.h
		0x41: packetClassUnassignChannel,
		0x42: packetClassAssignChannel,
		0x51: packetClassChannelID,
		0x43: packetClassChannelPeriod,
		0x44: packetClassSearchTKL,
		0x45: packetClassChannelFrequency,
		0x46: packetClassSetNetwork,
		0x47: packetClassTxPower,
		0x59: packetClassIdListAdd,
		0x5A: packetClassIdListConfig,
		0x60: packetClassChanneltxPower,
		0x63: packetClassLpSearchTimeout,
		0x65: packetClassSetSerialNumber,
		0x66: packetClassEnableExtMsgs,
		0x68: packetClassEnableLED,
		0x4A: packetClassSystemReset,
		0x4B: packetClassOpenChannel,
		0x4C: packetClassCloseChannel,
		0x5B: packetClassOpenRXScanCH,
		0x4D: packetClassReqMessage,
		0x4E: PacketClassBroadcastData,
		0x4F: packetClassAckData,
		0x50: packetClassBurstData,
		0x40: packetClassChannelEvent,
		0x52: packetClassChannelStatus,
		0x3E: packetClassVersion,
		0x54: packetClassCapabilities,
		0x61: packetClassSerialNumber,
		0x6F: packetClassNotifStartup,
		0x53: packetClassCwInit,
		0x48: packetClassCwTest,
		0xFF: packetClassUnknown,
	}
}

func PacketClassDecoder(classes map[byte]string) func(b byte) (string, error) {
	return func(b byte) (string, error) {
		class, ok := classes[b]
		if !ok {
			return packetClassUnknown, unknownPacketError(b)
		}
		return class, nil
	}
}

type unknownPacketError byte

func (dev unknownPacketError) Error() string {
	return fmt.Sprintf("UnknownHandler packet: %X", byte(dev))
}
