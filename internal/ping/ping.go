package ping

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/vldanch/dops/internal/config"
)

var (
	urls    []string
	timeout int
	retries int
)

// PingCmd represents the 'ping' command in the CLI
var PingCmd = &cobra.Command{
	Use:   "ping [--url URL]...",
	Short: "Ping one or more HTTP services concurrently",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if no URLs provided via flags or arguments
		if len(urls) == 0 && len(args) == 0 {
			fmt.Println("Please provide at least one URL using --url or as an argument")
			return
		}

		// If URLs are passed as positional arguments, add them to urls slice
		if len(args) > 0 {
			urls = append(urls, args...)
		}

		// Load configuration (to get default timeout and retries if not specified)
		cfg := config.Get()

		// Use config values if timeout or retries flags are not set
		if timeout == 0 {
			timeout = cfg.Ping.Timeout
		}
		if retries == 0 {
			retries = cfg.Ping.Retries
		}

		// WaitGroup to wait for all goroutines to finish
		var wg sync.WaitGroup
		wg.Add(len(urls))

		// Launch a goroutine per URL to ping concurrently
		for _, url := range urls {
			go func(u string) {
				defer wg.Done()
				pingService(u, timeout, retries)
			}(url)
		}

		// Wait until all pings are done
		wg.Wait()
	},
}

func init() {
	// Define CLI flags:
	// --url can be used multiple times to specify multiple URLs
	PingCmd.Flags().StringArrayVar(&urls, "url", []string{}, "URL(s) to ping, can be specified multiple times")
	// Timeout per ping attempt in seconds; overrides config if set
	PingCmd.Flags().IntVar(&timeout, "timeout", 0, "Timeout per ping attempt in seconds (overrides config)")
	// Number of retries per URL; overrides config if set
	PingCmd.Flags().IntVar(&retries, "retries", 0, "Number of retries per URL (overrides config)")
}

// pingService performs HTTP GET requests to the specified URL with timeout and retries
func pingService(url string, timeout, retries int) {
	// Create HTTP client with specified timeout
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	var lastErr error

	// Try to ping the URL up to 'retries' times
	for i := 1; i <= retries; i++ {
		start := time.Now()
		resp, err := client.Get(url)
		respTime := time.Since(start)

		if err != nil {
			// On error, save last error and print failure message
			lastErr = err
			fmt.Printf("Ping attempt %d to %s ... Failed: %v\n", i, url, err)
			continue
		}

		// On success, close response body and print success message with status code and response time
		resp.Body.Close()
		fmt.Printf("Ping attempt %d to %s ... Success! Status code: %d, Response time: %v\n", i, url, resp.StatusCode, respTime)
		return
	}

	// If all attempts failed, print a summary message with last error
	fmt.Printf("All %d ping attempts to %s failed. Last error: %v\n", retries, url, lastErr)
}
