package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vldanch/dops/cmd/checksystem"
	"github.com/vldanch/dops/cmd/ping"
	"github.com/vldanch/dops/pkg/config"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "dops",
	Short: "DevOps CLI Assistant",
	Long:  `dops - smart CLI tool for DevOps: ping, notifications, system checks.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use `dops --help` to see available commands.")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./configs/config.yaml)")

	rootCmd.AddCommand(ping.PingCmd)
	rootCmd.AddCommand(checksystem.Cmd)
}

func Execute() {
	cobra.OnInitialize(initConfig)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")
	}

	var cfg config.Config

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Config read error:", err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Println("Config unmarshal error:", err)
		os.Exit(1)
	}

	config.Set(cfg)
}
