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
	"q100paserver/current"
	"q100paserver/fan"
	"q100paserver/logger"
	"q100paserver/power"
	"q100paserver/temperature"
	"strings"
	"sync"
	"syscall"
)

const PORT = "9999"

func configureDevices() {
	power.Configure()
	current.Configure()
	temperature.Configure()
	fan.Configure()
}

// TODO: encode to json and include a version number (use json: tags).
// Could also have the client requst a version number to match
func readDevices() string {
	str := fmt.Sprintf("Pre %4.1f°, PA %4.1f° %3.1fA, Enc %04d->%04d, PA %04d->%04d",
		temperature.PreAmp(),
		temperature.FinalPA(),
		current.FinalPA(),
		fan.EnclosureIntake(),
		fan.EnclosureExtract(),
		fan.FinalPAintake(),
		fan.FinalPAextract(),
	)
	return str
}

func shutdownDevices() {
	power.Shutdown()
	logger.Info.Printf("Shutdown power        - done")
	fan.Shutdown()
	logger.Info.Printf("Shutdown fan          - done")
	current.Shutdown()
	logger.Info.Printf("Shutdown current      - done")
	temperature.Shutdown()
	logger.Info.Printf("Shutdown temperatutre - done")
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
		logger.Fatal.Fatalf("Failed to create listener: %s", err)
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
				logger.Warn.Println("accept error", err)
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

	logger.Info.Printf("got connection from: %v\n", conn.RemoteAddr())
	power.Up()
	clientReader := bufio.NewReader(conn)

	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')
		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if clientRequest == "CLOSE" {
				logger.Info.Printf("Connection closed with CLOSE")
				power.Down()
				return
			}
		case io.EOF:
			logger.Info.Printf("Connection closed with io.EOF")
			power.Down()
			return
		default:
			logger.Warn.Printf("Connection closed abnormally: %v", err)
			power.Down()
			return
		}

		// Responding to the client request (and check verion number match)
		str := readDevices() + "\n"

		if _, err = conn.Write([]byte(str)); err != nil {
			logger.Warn.Printf("failed to respond to client: %v\n", err)
		}
	}
}

func main() {
	logger.Open("/home/pi/Q100/paserver.log")
	defer logger.Close()

	logger.Info.Printf("Q-100 PA Server has started")

	// capture exit signals to ensure pin is reverted to input on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	configureDevices()

	s := NewServer("0.0.0.0:" + PORT)

	<-quit // wait on interupt

	logger.Info.Printf("---------- got interupt ----------")

	s.Stop()

	power.Down()
	shutdownDevices()
	logger.Info.Printf("Q-100 PA Server has stopped")

	// TODO: shutdown or reboot Rasberry Pi
}
