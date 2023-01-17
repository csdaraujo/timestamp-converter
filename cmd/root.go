/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "tc",
	Short:   "Simple CLI timestamp converter",
	Example: "tc 1593766111",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			// try to parse RFC3339 and convert to timestamp
			convertRFC3339ToUnixTimestamp(cmd, args[0])
		} else {
			convertUnixTimestampToRFC3339(cmd, args[0])
		}
	},
}

func convertRFC3339ToUnixTimestamp(cmd *cobra.Command, s string) {
	if strings.Contains(s, ".") { // has nano seconds
		parsed, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			cmd.PrintErrln("Invalid format!")
			return
		}
		cmd.Println(parsed.UnixNano())
	} else {
		parsed, err := time.Parse(time.RFC3339, s)
		if err != nil {
			cmd.PrintErrln("Invalid format!")
			return
		}
		cmd.Println(parsed.Unix())
	}
}

func convertUnixTimestampToRFC3339(cmd *cobra.Command, s string) {
	var parsedSeconds int64
	var parsedNanoseconds int64 = 0
	if len(s) > 10 {
		parsedSeconds, _ = strconv.ParseInt(s[:10], 10, 64)
		parsedNanoseconds, _ = strconv.ParseInt(s[10:], 10, 64)
	} else {
		parsedSeconds, _ = strconv.ParseInt(s, 10, 64)
	}

	tz := os.Getenv("TZ")
	locf := cmd.Flag("loc").Value.String()
	if locf != "" {
		tz = locf
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	converted := time.Unix(parsedSeconds, parsedNanoseconds).In(loc)
	if parsedNanoseconds == 0 {
		cmd.Println(converted.Format(time.RFC3339))
	} else {
		cmd.Println(converted.Format(time.RFC3339Nano))
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringP("loc", "l", "", "Set timezone")
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
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".timestamp-converter" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".timestamp-converter")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
