package ping

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/vldanch/dops/internal/config"
)

var (
	url     string
	timeout int
	retries int
)

var PingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping an HTTP service",
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			fmt.Println("Please provide a --url")
			return
		}

        // you can take values from the configuration if the flags are not set (for example, 0 or -1)
		cfg := config.Get()

        // If timeout is not specified in flags, we use it from the config
		if timeout <= 0 {
			timeout = cfg.Ping.Timeout
		}
		// If retries is not set, we use from config
		if retries <= 0 {
			retries = cfg.Ping.Retries
		}

		pingService(url, timeout, retries)
	},
}

func init() {
	PingCmd.Flags().StringVar(&url, "url", "", "URL to ping")
	PingCmd.Flags().IntVar(&timeout, "timeout", 0, "Timeout in seconds for each ping attempt")
	PingCmd.Flags().IntVar(&retries, "retries", 0, "Number of ping retries")
}

func pingService(url string, timeout, retries int) {
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	for i := 1; i <= retries; i++ {
		fmt.Printf("Ping attempt %d to %s ... ", i, url)
		start := time.Now()
		resp, err := client.Get(url)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Printf("Failed: %s\n", err)
			continue
		}

		fmt.Printf("Success! Status code: %d, Response time: %v\n", resp.StatusCode, elapsed)
		resp.Body.Close()
		return
	}
	fmt.Println("All ping attempts failed.")
}
