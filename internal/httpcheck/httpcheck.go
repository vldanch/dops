package httpcheck

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

// HttpCheckCmd represents the 'httpcheck' command in the CLI
var HttpCheckCmd = &cobra.Command{
	Use:     "httpcheck [--url URL]...",
	Aliases: []string{},
	Short:   "Check HTTP availability of one or more services concurrently",
	Run: func(cmd *cobra.Command, args []string) {
		if len(urls) == 0 && len(args) == 0 {
			fmt.Println("Please provide at least one URL using --url or as an argument")
			return
		}

		if len(args) > 0 {
			urls = append(urls, args...)
		}

		cfg := config.Get()

		if timeout == 0 {
			timeout = cfg.Ping.Timeout
		}
		if retries == 0 {
			retries = cfg.Ping.Retries
		}

		var wg sync.WaitGroup
		wg.Add(len(urls))

		for _, url := range urls {
			go func(u string) {
				defer wg.Done()
				checkHTTP(u, timeout, retries)
			}(url)
		}

		wg.Wait()
	},
}

func init() {
	HttpCheckCmd.Flags().StringArrayVar(&urls, "url", []string{}, "URL(s) to check, can be specified multiple times")
	HttpCheckCmd.Flags().IntVar(&timeout, "timeout", 0, "Timeout per HTTP request in seconds (overrides config)")
	HttpCheckCmd.Flags().IntVar(&retries, "retries", 0, "Number of retries per URL (overrides config)")
}

func checkHTTP(url string, timeout, retries int) {
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	var lastErr error

	for i := 1; i <= retries; i++ {
		start := time.Now()
		resp, err := client.Get(url)
		respTime := time.Since(start)

		if err != nil {
			lastErr = err
			fmt.Printf("HTTP attempt %d to %s ... Failed: %v\n", i, url, err)
			continue
		}

		resp.Body.Close()
		fmt.Printf("HTTP attempt %d to %s ... Success! Status code: %d, Response time: %v\n", i, url, resp.StatusCode, respTime)
		return
	}

	fmt.Printf("All %d HTTP attempts to %s failed. Last error: %v\n", retries, url, lastErr)
}
