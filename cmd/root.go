package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

type Config struct {
	Ping struct {
		Timeout int
		Retries int
	}
	Notify struct {
		TelegramToken string `mapstructure:"telegram_token"`
		ChatID        string `mapstructure:"chat_id"`
	}
}

var AppConfig Config

var rootCmd = &cobra.Command{
	Use:   "dops",
	Short: "DevOps CLI Assistant",
	Long:  `dops - smart CLI tool for DevOps: ping, notifications, system checks.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use `dops --help` to see available commands.")
	},
}

func Execute() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./configs/config.yaml)")
	rootCmd.AddCommand(PingCmd)

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
		viper.AddConfigPath("./configs") // in case config in configs/
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Config read error:", err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		fmt.Println("Config unmarshal error:", err)
		os.Exit(1)
	}

	SetConfig(AppConfig)
}
