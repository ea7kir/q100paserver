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
	"q100paserver/current"
	"q100paserver/fan"
	"q100paserver/logger"
	"q100paserver/power"
	"q100paserver/temperature"
	"strings"
)

const PORT = ":9999" // same as "0.0.0.0:9999"

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

// func prev_handler() {
// 	connected := true
// 	power.Up()
// 	for {
// 		if !connected {
// 			break
// 		}
// 		str := readDevices()
// 		logger.Info("\n\tSEND: %v\n", str)
// 		time.Sleep(2 * time.Second)
// 	}
// 	power.Down()
// }

// TODO: add signal to cancel

// http://www.inanzzz.com/index.php/post/j3n1/creating-a-concurrent-tcp-client-and-server-example-with-golang
func handleClientRequest(con net.Conn) {
	defer con.Close()

	logger.Info.Printf("got connection from: %v\n", con.RemoteAddr())
	power.Up()
	clientReader := bufio.NewReader(con)

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

		if _, err = con.Write([]byte(str)); err != nil {
			logger.Warn.Printf("failed to respond to client: %v\n", err)
		}
	}
}

// http://www.inanzzz.com/index.php/post/j3n1/creating-a-concurrent-tcp-client-and-server-example-with-golang

// 1st lets try the simple way

func runServer() {
	// capture exit signals to ensure pin is reverted to input on exit.
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// defer signal.Stop(quit)

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		logger.Fatal.Fatalf("Failed to create listener: %v", err) // TODO: sort out Fatal
	}
	defer listener.Close()

	logger.Info.Printf("Q-100 PA Server has started on %v", listener.Addr().String())

	for {
		// select {
		// case <-quit:
		// 	logger.Info.Printf("---------- got interupt ----------")
		// 	return
		// default:
		// }
		con, err := listener.Accept()
		if err != nil {
			logger.Warn.Printf("Accept failed: %v\n", err)
			continue
		}
		// If you want, you can increment a counter here and inject to handleClientRequest below as client identifier
		go handleClientRequest(con)
	}
}

// shutdown: https://eli.thegreenplace.net/2020/graceful-shutdown-of-a-tcp-server-in-go/

func main() {
	logger.Info.Printf("Q-100 PA Server has started")

	configureDevices()
	defer shutdownDevices()

	runServer()

	// shutdownDevices()
	logger.Info.Printf("Q-100 PA Server has stopped")
	// TODO: shutdown or reboot Rasberry Pi
}
