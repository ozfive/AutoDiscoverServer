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

type DiscoveryRequest struct {
	Command string
}

type DiscoveryResponse struct {
	Address string
}
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

	buffer := make([]byte, maxDatagramSize)

	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		var request DiscoveryRequest
		if err := decoder.Decode(&request); err != nil {
			fmt.Println("Error decoding request:", err)
			continue
		}

		if request.Command == "REQUEST" {
			localIP := GetLocalIP()
			response := DiscoveryResponse{Address: localIP + ":12345"}

			if err := encoder.Encode(response); err != nil {
				fmt.Println("Error encoding response:", err)
				continue
			}
			fmt.Println("Sent response to discovery request from", request.Command)
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
