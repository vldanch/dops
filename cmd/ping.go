package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var urlToPing string

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Pinging HTTP service",
	Run: func(cmd *cobra.Command, args []string) {
		if urlToPing == "" {
			fmt.Println("Please indicate the URL using -Url")
			return
		}
		client := http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Get(urlToPing)
		if err != nil {
			fmt.Printf("Request error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		fmt.Printf("Answer from %s: %d %s\n", urlToPing, resp.StatusCode, resp.Status)
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().StringVar(&urlToPing, "url", "", "URL for ping (for example, https://example.com)")
}
