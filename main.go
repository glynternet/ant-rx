package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/glynternet/ant-rx/antrx"
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

	app := antrx.App{
		PrintUnknown: *printUnknown,
		DebugMode:    *debug,
		DetectDevice: *detectDevice,
	}
	if err := app.Run(ctx); err != nil {
		fmt.Printf("error whilst running: %+v", err)
	}
}
