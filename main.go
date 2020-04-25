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

	"github.com/glynternet/ant-rx/ant"
	"github.com/glynternet/ant-rx/usbprompt"
	"github.com/google/gousb"
	"github.com/half2me/antgo/message"
	"github.com/pkg/errors"
)

func main() {
	debug := flag.Bool("debug", false, "debug logging")
	printUnknown := flag.Bool("print-unknown", false, "print unknown message types")
	detectDevice := flag.Bool("detect-device", false, "automatically detect ANT USB device")
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

	app := antRXApp{
		printUnknown: *printUnknown,
		debugMode:    *debug,
		detectDevice: *detectDevice,
	}
	if err := app.run(ctx); err != nil {
		fmt.Printf("error whilst running: %+v", err)
	}
}

type antRXApp struct {
	printUnknown bool
	debugMode    bool
	detectDevice bool
}

func (app antRXApp) run(ctx context.Context) error {
	usbCtx := gousb.NewContext()
	usbCtx.Debug(3)
	defer func() {
		if cErr := usbCtx.Close(); cErr != nil {
			log.Printf("Error closing USB context: %v", cErr)
		}
		if app.debugMode {
			log.Println("USB context closed")
		}
	}()

	device, err := app.usbDevice(usbCtx)
	if err != nil {
		return errors.Wrap(err, "getting ANT USB")
	}

	itf, close, err := setupInterface(app.debugMode, usbCtx, device)
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
	if err := sendRxScanModeMessages(sendCtx, app.debugMode, outep); err != nil {
		return errors.Wrap(err, "sending rx scan mode messages")
	}

	readstr, err := inep.NewStream(64, 1)
	if err != nil {
		return errors.Wrap(err, "preparing new in-endpoint stream for reading")
	}
	defer deferredClose(readstr, app.debugMode, "in-endpoint stream")

	fmt.Println("Listening for messages...")
	return handleMessages(ctx, readstr, ant.NewPacketHandler(newPacketPrinter(app.printUnknown)))
}

func (app antRXApp) usbDevice(usbCtx *gousb.Context) (*usbprompt.DeviceDesc, error) {
	if app.detectDevice {
		device, err := usbprompt.DetectDevice(usbCtx)
		if device != nil {
			fmt.Printf("Detected device: %s\n", device.HumanReadable())
		}
		return device, errors.Wrap(err, "getting detecting AND USB device")
	}
	device, err := usbprompt.UserSelectDevice(usbCtx)
	return device, errors.Wrap(err, "getting user selected USB device")
}

func setupInterface(debug bool, usbCtx *gousb.Context, device *usbprompt.DeviceDesc) (*gousb.Interface, func(), error) {
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

func handleMessages(ctx context.Context, str *gousb.ReadStream, handlePacket ant.PacketHandler) error {
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

			if buf[0] != message.MESSAGE_TX_SYNC {
				fmt.Printf("First byte was not ANT serial message Tx sync byte\n")
				continue
			}

			if err := handlePacket(buf); err != nil {
				return errors.Wrap(err, "handling ANT packet")
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
