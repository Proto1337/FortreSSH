package main

import (
    "fmt"
    "math/rand"
    "net"
    "strings"
    "time"
)

const (
    LISTEN_PORT = "2222"
)

func main() {
    // Start listening on the specified port
    listener, err := net.Listen("tcp", "0.0.0.0:" + LISTEN_PORT)
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
        go handleConnection(conn)
    }
    
}


func handleConnection(conn net.Conn) {
    // Handle the new connection, and write a long (and in fact, never-ending) banner
    
    // Close the connection when this function ends
    defer conn.Close()
    
    // We're going to be generating psuedo-random numbers, so seed it with the time the connection opened
    rand.Seed(time.Now().Unix())
    
    // Main loop - get a random string, write it, sleep then do it again
    for {
        var randstr = genString(10)
        conn.Write([]byte(randstr))
        conn.Write([]byte("\r\n"))
        // Sleep for 2s before sending the next
        time.Sleep(2 * time.Second)
    }
}


func genString(length int) (string){
    // Generate a psuedo-random string
    
    
    // Keep out charset ascii - it is pretending to be a printable banner, after all
    charSet := "abcdedfghijklmnopqrstABCDEDFGHIJKLMNOPQRSTUVWXYZ0123456789=.<>?!#@''"
    
    var output strings.Builder
    
    // Generate the string
    for i := 0; i < length; i++ {
        randnum := rand.Intn(len(charSet))
        randChar := charSet[randnum]
        output.WriteString(string(randChar))
    }
    
    return output.String()
    
}
