package cmd

import (
	"fmt"
	"log"

	"github.com/mparvin/tfd/tools"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tfd",
	Short: "Download file from Telegram",
	Long: `This bot will download the following content from Telegram:
    Video
    Documents
    Audios
    Links
    `,
	Run: func(cmd *cobra.Command, args []string) {
		// Ensure all directories exist
		if err := tools.EnsureDirectories(); err != nil {
			log.Fatalf("Failed to ensure directories: %v", err)
		}

		// Start the Telegram bot
		log.Println("Starting TFD (Telegram File Downloader)...")
		if err := tools.StartBot(); err != nil {
			log.Fatalf("Failed to start bot: %v", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tfd.yaml)")
	rootCmd.PersistentFlags().String("proxy", "", "proxy URL (optional, e.g., socks5://127.0.0.1:1080 or http://proxy.example.com:8080)")

	// Bind proxy flag to viper
	viper.BindPFlag("proxy", rootCmd.PersistentFlags().Lookup("proxy"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in home directory with name ".tfd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tfd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
