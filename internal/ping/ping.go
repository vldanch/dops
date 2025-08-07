package ping

import (
	"fmt"
	"net"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

func NewPingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ping [host]...",
		Short: "Ping one or more hosts to check network connectivity",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var hasErrors bool
			for _, host := range args {
				err := pingHost(host)
				if err != nil {
					fmt.Printf("Error pinging %s: %v\n", host, err)
					hasErrors = true
				}
				fmt.Println()
			}
			if hasErrors {
				return fmt.Errorf("some hosts could not be pinged")
			}
			return nil
		},
	}
	return cmd
}

func pingHost(host string) error {
	fmt.Printf("Pinging %s...\n", host)

	// TCP ping to port 80
	start := time.Now()
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, "80"), 2*time.Second)
	if err == nil {
		defer conn.Close()
		elapsed := time.Since(start).Milliseconds()
		fmt.Printf("TCP ping to %s: connected successfully in %d ms\n", host, elapsed)
		return nil
	}
	fmt.Printf("TCP ping failed: %v\n", err)

	// System ping fallback
	cmd := exec.Command("ping", "-c", "3", "-W", "2", host)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("System ping failed: %v\n", err)
		return fmt.Errorf("host unreachable by TCP and system ping")
	}

	fmt.Println(string(output))
	return nil
}
