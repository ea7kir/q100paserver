/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"

	"q100paserver/ds18b20monitor"
	"q100paserver/fandriver"
	"q100paserver/ina226monitor"
	"q100paserver/mylogger"
	"q100paserver/psudriver"
	"q100paserver/rpidriver"
	"strings"
	"sync"
	"syscall"
)

const kServerAddress = "0.0.0.0:9999" // "0.0.0.0:8765"

func configureDevices() {
	psudriver.Configure()
	ina226monitor.Configure()
	ds18b20monitor.Configure()
	fandriver.Configure()
	rpidriver.Configure()
}

// TODO: encode to json and include a version number (use json: tags).
// Could also have the client requst a version number to match
func readDevices() string {
	str := fmt.Sprintf("Pre %4.1f°, PA %4.1f° %3.1fA, Enc %04d->%04d, PA %04d->%04d, Pi %4.1f°",
		ds18b20monitor.PreAmpTemperature(),
		ds18b20monitor.PaTemperature(),
		ina226monitor.PaCurrent(),
		fandriver.EnclosureIntake(),
		fandriver.EnclosureExtract(),
		fandriver.FinalPAintake(),
		fandriver.FinalPAextract(),
		rpidriver.Temperature(),
	)
	return str
}

func shutdownDevices() {
	psudriver.Shutdown()
	mylogger.Info.Printf("Shutdown psudriver     - done")
	fandriver.Shutdown()
	mylogger.Info.Printf("Shutdown fandriver     - done")
	ina226monitor.Shutdown()
	mylogger.Info.Printf("Shutdown ina226driver  - done")
	ds18b20monitor.Shutdown()
	mylogger.Info.Printf("Shutdown ds18b20driver - done")
	rpidriver.Shutdown()
	mylogger.Info.Printf("Shutdown rpidriver     - done")
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
		mylogger.Fatal.Fatalf("Failed to create listener: %s", err)
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
				mylogger.Warn.Println("accept error", err)
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

	mylogger.Info.Printf("got connection from: %v\n", conn.RemoteAddr())
	psudriver.Up()
	clientReader := bufio.NewReader(conn)

	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')
		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if clientRequest == "CLOSE" {
				mylogger.Info.Printf("Connection closed with CLOSE")
				psudriver.Down()
				return
			}
		case io.EOF:
			mylogger.Info.Printf("Connection closed with io.EOF")
			psudriver.Down()
			return
		default:
			mylogger.Warn.Printf("Connection closed abnormally: %v", err)
			psudriver.Down()
			return
		}

		// Responding to the client request (and check verion number match)
		str := readDevices() + "\n"

		if _, err = conn.Write([]byte(str)); err != nil {
			mylogger.Warn.Printf("failed to respond to client: %v\n", err)
		}
	}
}

func main() {
	mylogger.Open("/home/pi/Q100/paserver.log")
	defer mylogger.Close()

	mylogger.Info.Printf("Q-100 PA Server has started")

	// capture exit signals to ensure pins are reverted to input on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	server := NewServer(kServerAddress)

	configureDevices()

	<-quit // wait on interupt

	mylogger.Info.Printf("---------- got interupt ----------")

	server.Stop()

	shutdownDevices()
	mylogger.Info.Printf("Q-100 PA Server has stopped")

	// TODO: shutdown or reboot Rasberry Pi
}
