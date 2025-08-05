package ping

import (
	"fmt"
	"time"
	"net"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func NewPingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ping [host]",
		Short: "Ping a host to check network connectivity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			host := args[0]
			return pingHost(host)
		},
	}
	return cmd
}

func pingHost(host string) error {
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("failed to resolve host: %v", err)
	}

	ip := ips[0].String()
	fmt.Printf("PING %s (%s) 56(84) bytes of data.\n", host, ip)

	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return fmt.Errorf("ListenPacket error: %v", err)
	}
	defer c.Close()

	for seq := 1; seq <= 3; seq++ {
		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  seq,
				Data: []byte("HELLO-R-U-THERE"),
			},
		}
		msgBytes, err := msg.Marshal(nil)
		if err != nil {
			return fmt.Errorf("Marshal error: %v", err)
		}

		start := time.Now()
		_, err = c.WriteTo(msgBytes, &net.IPAddr{IP: net.ParseIP(ip)})
		if err != nil {
			return fmt.Errorf("WriteTo error: %v", err)
		}

		c.SetReadDeadline(time.Now().Add(2 * time.Second))

		reply := make([]byte, 1500)
		n, peer, err := c.ReadFrom(reply)
		if err != nil {
			fmt.Printf("Request timeout for icmp_seq %d\n", seq)
			continue
		}

		duration := time.Since(start).Seconds() * 1000 // ms

		parsedMsg, err := icmp.ParseMessage(1, reply[:n])
		if err != nil {
			return fmt.Errorf("ParseMessage error: %v", err)
		}

		switch parsedMsg.Type {
		case ipv4.ICMPTypeEchoReply:
			fmt.Printf("64 bytes from %s: icmp_seq=%d time=%.2f ms\n", peer, seq, duration)
		default:
			fmt.Printf("Got %+v from %v; want echo reply\n", parsedMsg, peer)
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
