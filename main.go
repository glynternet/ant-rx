package main

import (
	"context"
	"flag"
	"fmt"
	"io"
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
	debug := flag.Bool("debug", false, "debug logging")
	printUnknown := flag.Bool("print-unknown", false, "print unknown message types")
	flag.Parse()
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
		printUnknown: *printUnknown,
		debug:        *debug,
	}); err != nil {
		fmt.Printf("error whilst running: %+v", err)
	}
}

type config struct {
	printUnknown bool
	debug        bool
}

func run(ctx context.Context, cfg config) error {
	usbCtx := gousb.NewContext()
	usbCtx.Debug(3)
	defer func() {
		if cErr := usbCtx.Close(); cErr != nil {
			log.Printf("Error closing USB context: %v", cErr)
		}
		if cfg.debug {
			log.Println("USB context closed")
		}
	}()

	chosen, err := userSelectDevice(usbCtx)
	if err != nil {
		return err
	}

	itf, close, err := setupInterface(cfg.debug, usbCtx, chosen)
	if err != nil {
		return errors.Wrap(err, "setting up USB interface")
	}
	defer close()

	inep, err := itf.InEndpoint(1)
	if err != nil {
		return errors.Wrap(err, "preparing in-endpoint")
	}

	outep, err := itf.OutEndpoint(1)
	if err != nil {
		return errors.Wrap(err, "preparing out-endpoint")
	}
	sendCtx, _ := context.WithTimeout(ctx, 10*time.Second)
	if err := sendRxScanModeMessages(sendCtx, cfg.debug, outep); err != nil {
		return errors.Wrap(err, "sending rx scan mode messages")
	}

	readstr, err := inep.NewStream(64, 1)
	if err != nil {
		return errors.Wrap(err, "preparing new in-endpoint stream for reading")
	}
	defer deferredClose(readstr, cfg.debug, "in-endpoint stream")

	p := printer{printUnknown: cfg.printUnknown}
	fmt.Println("Listening for messages...")
	return handleMessages(ctx, readstr, p)
}

func setupInterface(debug bool, usbCtx *gousb.Context, device *deviceDesc) (*gousb.Interface, func(), error) {
	dev, err := usbCtx.OpenDeviceWithVIDPID(device.Vendor, device.Product)
	if err != nil {
		return nil, nil, errors.Wrap(err, "opening device")
	}
	if dev == nil {
		return nil, nil, errors.New("open device is nil")
	}
	closeFns := []func(){deferredClose(dev, debug, "device")}
	if err := dev.SetAutoDetach(true); err != nil {
		for _, closeFn := range closeFns {
			closeFn()
		}
		return nil, nil, errors.Wrap(err, "setting up autodetach")
	}

	devCfg, err := dev.Config(1)
	if err != nil {
		for _, closeFn := range closeFns {
			closeFn()
		}
		return nil, nil, errors.Wrap(err, "something with configuration?!")
	}
	closeFns = append([]func(){deferredClose(devCfg, debug, "config")}, closeFns...)

	intf, err := devCfg.Interface(0, 0)
	if err != nil {
		for _, closeFn := range closeFns {
			closeFn()
		}
		return nil, nil, errors.Wrap(err, "claiming interface")
	}
	closeFns = append([]func(){func() {
		intf.Close()
		if debug {
			log.Println("interface closed")
		}
	}}, closeFns...)
	return intf, func() {
		for _, closeFn := range closeFns {
			closeFn()
		}
	}, nil
}

func handleMessages(ctx context.Context, str *gousb.ReadStream, v AntBroadcastMessageVisitor) error {
	classes := messageClasses()
	buf := make([]byte, 64)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			n, err := str.ReadContext(ctx, buf)
			if err != nil {
				log.Printf("Error reading from stream: %+v", err)
				return errors.Wrap(err, "reading from stream")
			}

			if n < 12 {
				continue
			}

			if buf[0] == message.MESSAGE_TX_SYNC {
				packet := message.AntPacket(buf)
				t, ok := classes[packet.Class()]
				if !ok {
					fmt.Printf("Unknown packet class: %d\n", packet.Class())
					continue
				}
				switch t {
				case messageClassBroadcastData:
					if err := VisitMessage(v, message.AntBroadcastMessage(packet)); err != nil {
						return errors.Wrap(err, "visiting message,")
					}
				default:
					fmt.Printf("Received packet: %s\n", t)
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

func deferredClose(c io.Closer, debug bool, name string) func() {
	return func() {
		if cErr := c.Close(); cErr != nil {
			log.Printf("Error closing Closer %s: %v", name, cErr)
		}
		if debug {
			log.Printf("device closed: %s", name)
		}
	}
}
