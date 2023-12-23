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
    "flag"
    "fmt"
    "math/rand"
    "net"
    "strings"
    "time"
)

var (
    LISTEN_PORT string  // Default listening port
    MIN_SLEEP   int     // Minimum time to respond between iterations
    MAX_SLEEP   int     // Maximum time to respond between iterations. Don't set too high or the client will timeout, suggest < 30
    MIN_LENGTH  int     // Minimum length of banner
    MAX_LENGTH  int     // This must not be set higher than 253 - SSH spec says 255 max, including the CRLF
)

func main() {
    /**
    * Start the listener on specified port.
    * Accept incoming connections.
    * Handle them in a goroutine.
    */
    // set the variables based on flags
    // no flags -> default values
    setVariables()
    // Start listening on the specified port
    listener, err := net.Listen("tcp", ":"+LISTEN_PORT)
    if err != nil {
        panic(err)
    }
    // Close the listener when the app closes
    defer listener.Close()
    fmt.Println("Listening on port " + LISTEN_PORT)
    // Start accepting connections in an endless loop
    // Print the incoming IP and send them into the goroutine
    for {
        conn, err := listener.Accept()
        if err != nil {
            // trap error and move onto a new connection
            continue
        }
        // Handle the connection in a new goroutine
        fmt.Println(time.Now().Format("2006/01/02-15:04:05"), conn.RemoteAddr().String(), "connected")
        go handleConnection(conn)
    }
}

func setVariables() {
    /*
    * Set the values of the variables based on commandline flags
    */
    flag.StringVar(&LISTEN_PORT, "port", "2222", "Port to listen on.")
    flag.IntVar(&MIN_SLEEP, "minsleep", 1, "Minimum time between iterations. Do not set it too high or the client will timeout, suggest < 30.")
    flag.IntVar(&MAX_SLEEP, "maxsleep", 5, "Maximum time between iterations. Do not set it too high or the client will timeout, suggest < 30.")
    flag.IntVar(&MIN_LENGTH, "minlength", 10, "Minimum length of responses. This must not be set higher than 253. The SSH specification says 255 max, including the CRLF.")
    flag.IntVar(&MAX_LENGTH, "maxlength", 120, "Maximum length of responses. This must not be set higher than 253. The SSH specification says 255 max, including the CRLF.")
    flag.Parse()
}

func handleConnection(conn net.Conn) {
    /**
    * Handle the new connection, and write a long (and in fact, never-ending) banner.
    */
    // Close the connection when this function ends
    defer conn.Close()
    // We're going to be generating psuedo-random numbers, so seed it with the time the connection opened
    start := time.Now().Unix()
    r := rand.New(rand.NewSource(start))
    // Main loop
    // Get a random string, write it, sleep then do it again
    for {
        // Generate random string and write it to the socket.
        _, err := conn.Write([]byte(genString(r.Intn(MAX_LENGTH-MIN_LENGTH)+MIN_LENGTH) + "\r\n"))
        // Now check the write worked - if the client went away we'll get an error
        // at that point, we should stop wasting resources and free up the FD
        if err != nil {
            now := time.Now()
            // Print length of connection
            fmt.Println(now.Format("2006/01/02-15:04:05"), conn.RemoteAddr().String(), "disconnected after", int(now.Unix()-start), "seconds")
            conn.Close()
            break
        }
        // Sleep for a period before sending the next
        // We vary the period a bit to tie the client up for varying amounts of time
        time.Sleep(time.Duration(r.Intn(MAX_SLEEP-MIN_SLEEP)+MIN_SLEEP) * time.Second)
    }
}

func genString(length int) string {
    /**
    * Generate a psuedo-random string
    * Keep out charset ascii - it is pretending to be a printable banner, after all
    */
    charSet := "abcdedfghijklmnopqrstABCDEDFGHIJKLMNOPQRSTUVWXYZ0123456789=.<>?!#@''"
    // Generate the string to send to the client
    var output strings.Builder
    for i := 0; i < length; i++ {
        randChar := charSet[rand.Intn(len(charSet))]
        output.WriteString(string(randChar))
    }
    return output.String()
}
