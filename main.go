package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	maxDatagramSize = 8192
	servicePort     = 10001
	requestString   = "REQUEST"
)

func main() {
	serverAddr := fmt.Sprintf(":%d", servicePort)

	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Handle program termination gracefully
	setupSignalHandler(conn)

	localIP := GetLocalIP()
	response := fmt.Sprintf("THE TCP SERVER IS LOCATED AT %s:12345", localIP)
	responseBytes := []byte(response)

	buffer := make([]byte, maxDatagramSize)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Println("Received", message, "from", addr)

		if message == requestString {
			_, err = conn.WriteToUDP(responseBytes, addr)
			if err != nil {
				fmt.Println("Error sending response:", err)
				continue
			}
			fmt.Println("Sent", response, "to", addr)
		}
	}
}

// GetLocalIP returns the non-loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// setupSignalHandler sets up a signal handler to clean up resources on program termination
func setupSignalHandler(conn *net.UDPConn) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nReceived termination signal. Cleaning up resources...")
		conn.Close()
		os.Exit(0)
	}()
}
