package main

const (
	messageClassUnassignChannel  = "unassign_channel"
	messageClassAssignChannel    = "assign_channel"
	messageClassChannelID        = "channel_id"
	messageClassChannelPeriod    = "channel_period"
	messageClassSearchTKL        = "search_timeout"
	messageClassChannelFrequency = "channel_frequency"
	messageClassSetNetwork       = "set_network"
	messageClassTxPower          = "tx_power"
	messageClassIdListAdd        = "id_list_add"
	messageClassIdListConfig     = "id_list_config"
	messageClassChanneltxPower   = "channel_tx_power"
	messageClassLpSearchTimeout  = "lp_search_timeout"
	messageClassSetSerialNumber  = "set_serial_number"
	messageClassEnableExtMsgs    = "enable_ext_msgs"
	messageClassEnableLED        = "enable_led"
	messageClassSystemReset      = "system_reset"
	messageClassOpenChannel      = "open_channel"
	messageClassCloseChannel     = "close_channel"
	messageClassOpenRXScanCH     = "open_rx_scan_ch"
	messageClassReqMessage       = "req_message"
	messageClassBroadcastData    = "broadcast_data"
	messageClassAckData          = "ack_data"
	messageClassBurstData        = "burst_data"
	messageClassChannelEvent     = "channel_event"
	messageClassChannelStatus    = "channel_status"
	messageClassVersion          = "version"
	messageClassCapabilities     = "capabilities"
	messageClassSerialNumber     = "serial_number"
	messageClassNotifStartup     = "notif_startup"
	messageClassCwInit           = "cw_init"
	messageClassCwTest           = "cw_test"
	messageClassUnknown          = "unhandled"
)

func packetClasses() map[byte]string {
	return map[byte]string{
		// From here: https://github.com/GoldenCheetah/GoldenCheetah/blob/3a31f5d131df46c90e25810a876ee4c5e0db5512/src/ANT/ANT.h
		0x41: messageClassUnassignChannel,
		0x42: messageClassAssignChannel,
		0x51: messageClassChannelID,
		0x43: messageClassChannelPeriod,
		0x44: messageClassSearchTKL,
		0x45: messageClassChannelFrequency,
		0x46: messageClassSetNetwork,
		0x47: messageClassTxPower,
		0x59: messageClassIdListAdd,
		0x5A: messageClassIdListConfig,
		0x60: messageClassChanneltxPower,
		0x63: messageClassLpSearchTimeout,
		0x65: messageClassSetSerialNumber,
		0x66: messageClassEnableExtMsgs,
		0x68: messageClassEnableLED,
		0x4A: messageClassSystemReset,
		0x4B: messageClassOpenChannel,
		0x4C: messageClassCloseChannel,
		0x5B: messageClassOpenRXScanCH,
		0x4D: messageClassReqMessage,
		0x4E: messageClassBroadcastData,
		0x4F: messageClassAckData,
		0x50: messageClassBurstData,
		0x40: messageClassChannelEvent,
		0x52: messageClassChannelStatus,
		0x3E: messageClassVersion,
		0x54: messageClassCapabilities,
		0x61: messageClassSerialNumber,
		0x6F: messageClassNotifStartup,
		0x53: messageClassCwInit,
		0x48: messageClassCwTest,
		0xFF: messageClassUnknown,
	}
}
