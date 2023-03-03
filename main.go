package main

import (
	"fmt"
	"net"
	"os"
)

const maxDatagramSize = 8192

func main() {

	/* Prepare a address at any address at port 10001*/
	ServerAddr, err := net.ResolveUDPAddr("udp", ":10001")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer ServerConn.Close()

	writedata := "THE TCP SERVER IS LOCATED AT " + GetLocalIP() + ":12345"

	buf := make([]byte, maxDatagramSize)
	writebuf := []byte(writedata)
	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		message := make([]byte, n)
		copy(message, buf[:n])
		fmt.Println("Received", string(message), " from ", addr)
		if string(message) == "REQUEST" {
			_, err = ServerConn.WriteToUDP(writebuf, addr)
			if err != nil {
				fmt.Println("Error: ", err)
				continue
			}
			fmt.Println("Sent", string(writebuf), "from", GetLocalIP()+":10001", "to", addr)
		}
	}
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type, if it is not a loopback address then display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
