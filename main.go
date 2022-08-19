package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	listen     = flag.String("l", ":1514", "listening address")
	bufferSize = flag.Int("b", 8192, "message buffer size")
	target     = flag.String("t", "127.0.0.1:1513", "target address")
)

func main() {
	buff := make([]byte, *bufferSize)
	srv, err := net.ListenPacket("udp", *listen)
	if err != nil {
		log.Fatalln(err)
	}
	raddr, err := net.ResolveUDPAddr("udp", *target)
	if err != nil {
		log.Fatalln(err)
	}
	r, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Listening on:", *listen)
	log.Println("Target:", *target)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
loop:
	for {
		select {
		case <-sigCh:
			break loop
		default:
		}
		var n int
		var err error
		if err := srv.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
			log.Fatalln(err)
		}
		n, _, err = srv.ReadFrom(buff)
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				continue
			}
			log.Fatalln(err)
		}
		buff[n] = '\n'
		n++
		r.Write(buff[:n])
	}
	log.Println("Bye")
}
