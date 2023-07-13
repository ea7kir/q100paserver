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
	"q100paserver/current"
	"q100paserver/fan"
	"q100paserver/logger"
	"q100paserver/power"
	"q100paserver/temperature"
	"strings"
)

func configureDevices() {
	pi := 0
	power.Configure(pi)
	current.Configure(pi)
	temperature.Configure(pi)
	fan.Configure(pi)
}

// func readDevices() string {
// 	str := fmt.Sprintf("%v %v %v",
// 		temperature.Read(),
// 		current.Read(),
// 		fan.Read())
// 	return str
// }

func shutdownDevices() {
	power.Shutdown()
	fan.Shutdown()
	current.Shutdown()
	temperature.Shutdown()

}

// func prev_handler() {
// 	connected := true
// 	power.Up()
// 	for {
// 		if !connected {
// 			break
// 		}
// 		str := readDevices()
// 		logger.Info.Printf("\n\tSEND: %v\n", str)
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
				// log.Println("client requested server to close the connection so closing")
				logger.Info.Printf("Connection closed with CLOSE")
				power.Down()
				return
			} else {
				// log.Printf("Received >%v<", clientRequest)
				logger.Info.Printf("Received >%v<", clientRequest)
			}
		case io.EOF:
			// log.Println("client closed the connection by terminating the process")
			logger.Info.Printf("Connection closed with io.EOF")
			power.Down()
			return
		default:
			// log.Printf("error: %v\n", err)
			logger.Warn.Printf("Connection closed abnormally: %v", err)
			power.Down()
			return
		}

		// Responding to the client request
		str := fmt.Sprintf("%v %v %v\n",
			temperature.Read(),
			current.Read(),
			fan.Read())
		// if _, err = con.Write([]byte("GOT IT!\n")); err != nil {
		if _, err = con.Write([]byte(str)); err != nil {
			// log.Printf("failed to respond to client: %v\n", err)
			logger.Warn.Printf("failed to respond to client: %v\n", err)
		}
	}
}

// http://www.inanzzz.com/index.php/post/j3n1/creating-a-concurrent-tcp-client-and-server-example-with-golang
func runServer() {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err) // TODO: sort out Fatal
	}
	defer listener.Close()

	for {
		con, err := listener.Accept()
		if err != nil {
			// log.Println(err)
			logger.Warn.Printf("Accept failed: %v\n", err)
			continue
		}

		// If you want, you can increment a counter here and inject to handleClientRequest below as client identifier
		go handleClientRequest(con)
	}
}

func main() {
	logger.Info.Printf("Q-100 PA Server will start...")
	configureDevices()
	runServer()
	shutdownDevices()
	logger.Info.Printf("Q-100 PA Server has stopped")
	// TODO: shutdown or reboot Rasberry Pi
}
