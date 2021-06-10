/*
Copyright © 2021 Bhargav SNV bhargavsnv100@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"golang.org/x/sync/errgroup"

	"github.com/Gituser143/cryptgo/pkg/api"
	"github.com/Gituser143/cryptgo/pkg/display/allcoin"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cryptgo",
	Short: "A terminal application to watch crypto prices!",
	Long:  `Crytpgo is a TUI based application written purely in Go to monitor and observe cryptocurrency prices in real time!`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Context and errgroup used to manage routines
		eg, ctx := errgroup.WithContext(context.Background())
		dataChannel := make(chan api.AssetData)
		searchChannel := make(chan []api.CoinSearchDetails)

		// Flag to determine if data must be sent when viewing per coin prices
		sendData := true

		// Fetch Coin Assets
		eg.Go(func() error {
			return api.GetAssets(ctx, dataChannel, &sendData)
		})

		// Fetch Top 3 coin history
		eg.Go(func() error {
			return api.GetTopCoinData(ctx, dataChannel, &sendData)
		})

		// Display UI for overall coins
		eg.Go(func() error {
			return allcoin.DisplayAllCoins(ctx, dataChannel, searchChannel, &sendData)
		})

		if err := eg.Wait(); err != nil {
			if err.Error() != "UI Closed" {
				return err
			}
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cryptgo.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cryptgo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cryptgo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
