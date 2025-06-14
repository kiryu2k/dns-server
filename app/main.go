package main

import (
	"fmt"
	"net"

	"github.com/kiryu2k/dns-server/internal/domain"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer func() {
		if err := udpConn.Close(); err != nil {
			fmt.Println("failed to close upd connection:", err)
		}
	}()

	buf := make([]byte, 512)
	for {
		_, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		msg, err := domain.MessageFromBytes(buf)
		if err != nil {
			fmt.Println("new message:", err)
			break
		}

		response := domain.NewMessage(msg.Header.Id).
			AsReply().
			WithQuestion(msg.Question).
			Encode()
		if _, err = udpConn.WriteToUDP(response, source); err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
