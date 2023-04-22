/*
Copyright © 2023 SWIFT_DEVS <https://github.com/321swift>
*/
package cmd

import (
	"os"

	"github.com/321swift/swift/client"
	"github.com/321swift/swift/server"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "swift",
	Short: "A tool for sending and receiving files to other pcs in your network",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		switch WelcomeScreen() {
		case 1:
			sender := server.NewServer()
			sender.Start()
		case 2:
			receiver := client.NewClient()
			receiver.Listen()
			// receiver.Connect()
		default:
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.swift.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
