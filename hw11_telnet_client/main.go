package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", time.Second*10, "connection timeout")
	flag.Parse()
	if len(flag.Args()) < 2 {
		log.Fatal("host or port not sent")
	}
	host := flag.Arg(0)
	port := flag.Arg(1)
	if _, err := strconv.Atoi(port); err != nil {
		fmt.Printf("%q must be a number.\n", port)
	}

	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	log.Printf("Start connection to %s", address)
	if err := client.Connect(); err != nil {
		log.Fatalln(err)
	}

	defer func(client TelnetClient) {
		if err := client.Close(); err != nil {
			log.Fatalln(err)
		}
	}(client)

	log.Printf("Connected to %s \n", address)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		<-ctx.Done()
		cancel()
	}()

	go func() {
		if err := client.Send(); err != nil {
			log.Fatalln(err)
		}
		log.Print("EOF")
		cancel()
	}()

	go func() {
		err := client.Receive()
		if err != nil {
			log.Fatalln(err)
		}
		log.Print("Connection was closed by peer")
		cancel()
	}()
}
