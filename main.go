package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/gousb"
	"github.com/half2me/antgo/message"
	"github.com/pkg/errors"
)

func main() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for range sigs {
			fmt.Printf("Received signal, cancelling...")
			cancel()
		}
	}()
	defer cancel()
	if err := run(ctx, config{
		debug:        true,
	}); err != nil {
		fmt.Printf("error whilst running: %+v", err)
	}
}

type config struct {
	debug        bool
}

func run(ctx context.Context, cfg config) error {
	usbCtx := gousb.NewContext()
	usbCtx.Debug(3)
	defer usbCtx.Close()

	chosen, err := userSelectDevice(usbCtx)
	if err != nil {
		return err
	}

	dev, err := usbCtx.OpenDeviceWithVIDPID(chosen.Vendor, chosen.Product)
	if err != nil {
		return errors.Wrap(err, "opening device")
	}
	if dev == nil {
		return errors.New("open device is nil")
	}
	defer func() {
		if cErr := dev.Close(); cErr != nil {
			log.Printf("Error closing device: %v", cErr)
		}
		if cfg.debug {
			log.Println("device closed")
		}
	}()

	if err := dev.SetAutoDetach(true); err != nil {
		return errors.Wrap(err, "setting up autodetach")
	}

	devCfg, err := dev.Config(1)
	if err != nil {
		return errors.Wrap(err, "something with configuration?!")
	}
	defer func() {
		if cErr := devCfg.Close(); cErr != nil {
			log.Printf("Error closing config: %v", cErr)
		}
		if cfg.debug {
			log.Println("config closed")
		}
	}()

	intf, err := devCfg.Interface(0, 0)
	if err != nil {
		return errors.Wrap(err, "claiming interface")
	}
	defer func() {
		intf.Close()
		if cfg.debug {
			log.Println("interface closed")
		}
	}()

	inep, err := intf.InEndpoint(1)
	if err != nil {
		return errors.Wrap(err, "preparing in-endpoint")
	}

	outep, err := intf.OutEndpoint(1)
	if err != nil {
		return errors.Wrap(err, "preparing out-endpoint")
	}

	readstr, err := inep.NewStream(64, 1)
	if err != nil {
		return errors.Wrap(err, "preparing new in-endpoint stream for reading")
	}
	defer func() {
		if cErr := readstr.Close(); cErr != nil {
			log.Printf("Error closing in-endpoint stream: %v", cErr)
		}
		if cfg.debug {
			log.Println("readstream closed")
		}
	}()

	// TODO(glynternet): can I do something with the main timeout here?
	sendCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := sendRxScanModeMessages(sendCtx, cfg.debug, outep); err != nil {
		return errors.Wrap(err, "sending rx scan mode messages")
	}

	opCtx := context.Background()
	buf := make([]byte, 64)
	var p printer
	fmt.Println("Listening for messages...")
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			n, err := readstr.ReadContext(opCtx, buf)
			if err != nil {
				log.Printf("Error reading from stream: %+v", err)
				return errors.Wrap(err, "reading from stream")
			}

			if n < 12 {
				continue
			}

			if buf[0] == message.MESSAGE_TX_SYNC {
				packet := message.AntPacket(buf)
				if packet.Class() == message.MESSAGE_TYPE_BROADCAST {
					if err := VisitMessage(p, message.AntBroadcastMessage(packet)); err != nil {
						return errors.Wrap(err, "visiting message")
					}
				}
			}
		}
	}
}

func sendRxScanModeMessages(ctx context.Context, debug bool, ep *gousb.OutEndpoint) error {
	for _, packet := range []struct {
		packet message.AntPacket
		desc   string
	}{{
		packet: message.SystemResetMessage(),
		desc:   "system reset message",
	}, {
		packet: message.SetNetworkKeyMessage(0, []byte(message.ANTPLUS_NETWORK_KEY)),
		desc:   "set network key message",
	}, {
		packet: message.AssignChannelMessage(0, message.CHANNEL_TYPE_ONEWAY_RECEIVE),
		desc:   "assign channel message",
	}, {
		packet: message.SetChannelIdMessage(0),
		desc:   "set channel ID message",
	}, {
		packet: message.SetChannelRfFrequencyMessage(0, 2457),
		desc:   "set channel RF Frequency message",
	}, {
		packet: message.EnableExtendedMessagesMessage(true),
		desc:   "enable extended messages message",
	}, {
		packet: message.LibConfigMessage(true, true, true),
		desc:   "lib config message",
	}, {
		packet: message.OpenRxScanModeMessage(),
		desc:   "open rx scan mode message",
	}} {
		if _, err := ep.WriteContext(ctx, packet.packet); err != nil {
			return errors.Wrapf(err, "sending message: %s", packet.desc)
		}
		if debug {
			fmt.Printf("Message sent: %s\n", packet.desc)
		}
	}
	return nil
}
