/**
* FortreSSH
*
* This program accepts incoming connections from SSH clients and
* send them an infinitly long SSH banner.
* The idea of so-called "tarpits" is to have automated scripts wasting
* resources on your tarpit instead of bothering others.
* 
* This program is a fork of Ben Tasker's 'Golang-SSH-Tarpit'.
* 
* Golang-SSH-Tarpit Copyright (C) 2021 B Tasker
* FortreSSH Copyright (C) 2023 Umut "proto" Yilmaz
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

const (
	LISTEN_PORT = "2222"    // Default listening port
	MIN_SLEEP   = 1         // Minimum time to respond between iterations
	MAX_SLEEP   = 5         // Maximum time to respond between iterations. Don't set too high or the client will timeout, suggest < 30
	MIN_LENGTH  = 10        // Minimum length of banner
	MAX_LENGTH  = 120       // This must not be set higher than 253 - SSH spec says 255 max, including the CRLF
)

func main() {
	// Start listening on the specified port
	listener, err := net.Listen("tcp", "0.0.0.0:"+LISTEN_PORT)
	if err != nil {
		panic(err)
	}

	// Close the listener when the app closes
	defer listener.Close()
	fmt.Println("Listening on " + LISTEN_PORT)

	// Start accepting connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			// trap error and move onto a new connection
			continue
		}

		// Handle the connection in a new goroutine
		now := time.Now().Format("2006/01/02-15:04:05")
		fmt.Println(now, "Tarpitting "+conn.RemoteAddr().String())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Handle the new connection, and write a long (and in fact, never-ending) banner

	// Close the connection when this function ends
	defer conn.Close()

	// We're going to be generating psuedo-random numbers, so seed it with the time the connection opened
	start := time.Now().Unix()
	r := rand.New(rand.NewSource(start))

	// Main loop - get a random string, write it, sleep then do it again
	for {
		// Generate random string and write it to the socket.
		_, err := conn.Write([]byte(genString(r.Intn(MAX_LENGTH-MIN_LENGTH)+MIN_LENGTH) + "\r\n"))

		// Now check the write worked - if the client went away we'll get an error
		// at that point, we should stop wasting resources and free up the FD
		if err != nil {
			now := time.Now()
			// Print length of connection
			fmt.Println(now.Format("2006/01/02-15:04:05"), "Coward disconnected:", conn.RemoteAddr().String(), "after", int(now.Unix()-start), "seconds")
			conn.Close()
			break
		}

		// Sleep for a period before sending the next
		// We vary the period a bit to tie the client up for varying amounts of time
		time.Sleep(time.Duration(r.Intn(MAX_SLEEP-MIN_SLEEP)+MIN_SLEEP) * time.Second)
	}
}

func genString(length int) string {
	// Generate a psuedo-random string
	// Keep out charset ascii - it is pretending to be a printable banner, after all
	charSet := "abcdedfghijklmnopqrstABCDEDFGHIJKLMNOPQRSTUVWXYZ0123456789=.<>?!#@''"

	var output strings.Builder

	// Generate the string
	for i := 0; i < length; i++ {
		randChar := charSet[rand.Intn(len(charSet))]
		output.WriteString(string(randChar))
	}

	return output.String()
}
