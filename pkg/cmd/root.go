// Copyright © 2019 WEBOP NINJAS, LLC <shawn@webop.ninja>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var awsKey string
var awsSecret string
var awsRegion string


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "eks-iam-group-mapper update",
	Short: "Maps IAM group members to kubernetes group members",
	Long: `Maps IAM group members to kubernetes group members`,
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&awsKey, "aws-access-key-id", "", "AWS Access Key ID")
	rootCmd.PersistentFlags().StringVar(&awsSecret, "aws-secret-access-key", "", "AWS Secret Access Key")
	rootCmd.PersistentFlags().StringVar(&awsRegion, "aws-region", "us-east-1", "AWS Region")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
