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

// Package cmd provides all the commands of the cli tool.
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// var path string.
var path string

// NewRootCmd returns the root command.
func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "seeder",
		Short: "Database seeds. ClI and Golang library",
		Long: `Seeder is a ClI tool and Golang library that helps to
seeds databases using golang code. ORM or SQL driver agnostic.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var out bytes.Buffer

			c := exec.Command("go", "run", fmt.Sprint(path, "/", "main.go"))
			c.Stdout = &out

			if err := c.Run(); err != nil {
				return err
			}

			return nil
		},
	}
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = NewRootCmd()

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.seeder.yaml)")
	dir, _ := os.Getwd()
	rootCmd.Flags().StringVarP(&path, "path", "p", filepath.Join(dir, "db"), "")
}

// initConfig reads in config file and ENV variables if set.
// func initConfig() {
// if cfgFile != "" {
// // Use config file from the flag.
// viper.SetConfigFile(cfgFile)
// } else {
// // Find home directory.
// home, err := os.UserHomeDir()
// cobra.CheckErr(err)
//
// // Search config in home directory with name ".seeder" (without extension).
// viper.AddConfigPath(home)
// viper.SetConfigType("yaml")
// viper.SetConfigName(".seeder")
// }
//
// viper.AutomaticEnv() // read in environment variables that match
//
// // If a config file is found, read it in.
// if err := viper.ReadInConfig(); err == nil {
// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
// }
// }
