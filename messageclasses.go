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
)

func messageClasses() map[byte]string {
	return map[byte]string{
		// From here: https://github.com/GoldenCheetah/GoldenCheetah/blob/3a31f5d131df46c90e25810a876ee4c5e0db5512/src/ANT/ANT.h
		0x41: messageClassUnassignChannel,  //ANT_UNASSIGN_CHANNEL
		0x42: messageClassAssignChannel,    //ANT_ASSIGN_CHANNEL
		0x51: messageClassChannelID,        //ANT_CHANNEL_ID
		0x43: messageClassChannelPeriod,    //ANT_CHANNEL_PERIOD
		0x44: messageClassSearchTKL,        //ANT_SEARCH_TIMEOUT
		0x45: messageClassChannelFrequency, //ANT_CHANNEL_FREQUENCY
		0x46: messageClassSetNetwork,       //ANT_SET_NETWORK
		0x47: messageClassTxPower,          //ANT_TX_POWER
		0x59: messageClassIdListAdd,        //ANT_ID_LIST_ADD
		0x5A: messageClassIdListConfig,     //ANT_ID_LIST_CONFIG
		0x60: messageClassChanneltxPower,   //ANT_CHANNEL_TX_POWER
		0x63: messageClassLpSearchTimeout,  //ANT_LP_SEARCH_TIMEOUT
		0x65: messageClassSetSerialNumber,  //ANT_SET_SERIAL_NUMBER
		0x66: messageClassEnableExtMsgs,    //ANT_ENABLE_EXT_MSGS
		0x68: messageClassEnableLED,        //ANT_ENABLE_LED
		0x4A: messageClassSystemReset,      //ANT_SYSTEM_RESET
		0x4B: messageClassOpenChannel,      //ANT_OPEN_CHANNEL
		0x4C: messageClassCloseChannel,     //ANT_CLOSE_CHANNEL
		0x5B: messageClassOpenRXScanCH,     //ANT_OPEN_RX_SCAN_CH
		0x4D: messageClassReqMessage,       //ANT_REQ_MESSAGE
		0x4E: messageClassBroadcastData,    //ANT_BROADCAST_DATA
		0x4F: messageClassAckData,          //ANT_ACK_DATA
		0x50: messageClassBurstData,        //ANT_BURST_DATA
		0x40: messageClassChannelEvent,     //ANT_CHANNEL_EVENT
		0x52: messageClassChannelStatus,    //ANT_CHANNEL_STATUS
		0x3E: messageClassVersion,          //ANT_VERSION
		0x54: messageClassCapabilities,     //ANT_CAPABILITIES
		0x61: messageClassSerialNumber,     //ANT_SERIAL_NUMBER
		0x6F: messageClassNotifStartup,     //ANT_NOTIF_STARTUP
		0x53: messageClassCwInit,           //ANT_CW_INIT
		0x48: messageClassCwTest,           //ANT_CW_TEST
	}
}
