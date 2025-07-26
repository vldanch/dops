package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var PingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping an HTTP service",
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			fmt.Println("Please provide a --url")
			return
		}
		pingService(url, config.Ping.Timeout, config.Ping.Retries)
	},
}

var url string
var config Config

func init() {
	PingCmd.Flags().StringVar(&url, "url", "", "URL to ping")
}

func SetConfig(c Config) {
	config = c
}

func pingService(url string, timeout, retries int) {
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	for i := 1; i <= retries; i++ {
		fmt.Printf("Ping attempt %d to %s ... ", i, url)
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("Failed: %s\n", err)
			continue
		}
		fmt.Printf("Success! Status code: %d\n", resp.StatusCode)
		resp.Body.Close()
		return
	}
	fmt.Println("All ping attempts failed.")
}
