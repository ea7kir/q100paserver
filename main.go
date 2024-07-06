/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"

	"q100paserver/ds18b20monitor"
	"q100paserver/fanMonitor"
	"q100paserver/ina226monitor"
	"q100paserver/psuSwitcher"
	"q100paserver/rpiMonitor"
	"strings"
	"sync"
	"syscall"
)

const kServerAddress = "0.0.0.0:9999" // "0.0.0.0:8765"

func configureDevices() {
	psuSwitcher.Configure()
	ina226monitor.Configure()
	ds18b20monitor.Configure()
	fanMonitor.Configure()
	rpiMonitor.Configure()
}

// TODO: encode to json and include a version number (use json: tags).
// Could also have the client requst a version number to match
func readDevices() string {
	str := fmt.Sprintf("Pre %4.1f°, PA %4.1f° %3.1fA, Enc %04d->%04d, PA %04d->%04d, Pi %4.1f°",
		ds18b20monitor.PreAmpTemperature(),
		ds18b20monitor.PaTemperature(),
		ina226monitor.PaCurrent(),
		fanMonitor.EnclosureIntake(),
		fanMonitor.EnclosureExtract(),
		fanMonitor.FinalPAintake(),
		fanMonitor.FinalPAextract(),
		rpiMonitor.Temperature(),
	)
	return str
}

func shutdownDevices() {
	psuSwitcher.Shutdown()
	log.Printf("INFO Shutdown psudriver     - done")
	fanMonitor.Shutdown()
	log.Printf("INFO Shutdown fandriver     - done")
	ina226monitor.Shutdown()
	log.Printf("INFO Shutdown ina226driver  - done")
	ds18b20monitor.Shutdown()
	log.Printf("INFO Shutdown ds18b20driver - done")
	rpiMonitor.Shutdown()
	log.Printf("INFO Shutdown rpidriver     - done")
}

// https://eli.thegreenplace.net/2020/graceful-shutdown-of-a-tcp-server-in-go/
// eg: https://github.com/eliben/code-for-blog/blob/master/2020/tcp-server-shutdown/shutdown1/shutdown1.go

// Socket server that can be shut down -- stop serving, in a graceful manner.
// This version expects all clients to close their connections before it
// successfully returns from Stop().
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.

type Server struct {
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
}

func NewServer(addr string) *Server {
	s := &Server{
		quit: make(chan interface{}),
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("FATAL Failed to create listener: %s", err)
	}
	s.listener = l
	s.wg.Add(1)
	go s.serve()
	return s
}

func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}

func (s *Server) serve() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Printf("WARN accept error: %s", err)
			}
		} else {
			s.wg.Add(1)
			go func() {
				s.handleConection(conn)
				s.wg.Done()
			}()
		}
	}
}

func (s *Server) handleConection(conn net.Conn) {
	defer conn.Close()

	log.Printf("INFO got connection from: %v\n", conn.RemoteAddr())
	psuSwitcher.Up()
	clientReader := bufio.NewReader(conn)

	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')
		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if clientRequest == "CLOSE" {
				log.Printf("INFO Connection closed with CLOSE")
				psuSwitcher.Down()
				return
			}
		case io.EOF:
			log.Printf("INFO Connection closed with io.EOF")
			psuSwitcher.Down()
			return
		default:
			log.Printf("WARN Connection closed abnormally: %s", err)
			psuSwitcher.Down()
			return
		}

		// Responding to the client request (and check verion number match)
		str := readDevices() + "\n"

		if _, err = conn.Write([]byte(str)); err != nil {
			log.Printf("WARN failed to respond to client: %s\n", err)
		}
	}
}

func main() {

	log.Printf("INFO Q-100 PA Server has started")

	// capture exit signals to ensure pins are reverted to input on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	server := NewServer(kServerAddress)

	configureDevices()

	<-quit // wait on interupt

	log.Printf("INFO ---------- got interupt ----------")

	server.Stop()

	shutdownDevices()
	log.Printf("INFO Q-100 PA Server has stopped")

	// TODO: shutdown or reboot Rasberry Pi
}
