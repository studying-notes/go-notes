/*
 * @Date: 2022.01.12 10:34
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.12 10:34
 */

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 写法让人难受，不喜欢，不用

var (
	// Used for flags.
	configFIle  string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "cobra",
		Short: "A generator for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of cobra",
		Long:  `All software has versions. This is cobra's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Cobra v0.9 -- HEAD")
		},
	}

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts web server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starts web server")
		},
	}

	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Show all config",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%+v", viper.AllSettings())
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFIle, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(configCmd)
}

func initConfig() {
	if configFIle != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFIle)
	} else {
		// Find home directory.
		home, _ := os.UserHomeDir()

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	Execute()
}
