/*
Copyright © 2020 David Arnold <dar@xoe.solutions>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/xoe-labs/go-generators/ddd-domain-gen/pkg/generate"
)

var (
	cfgFile    string
	sourceType string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ddd-domain-gen",
	Short: "Generates idiomatic go code for a DDD domain",
	Long: `Generates idiomatic go code for a DDD domain based on struct field annotations.

  Available Annotations:
    gen-getter               - generate default getter: no special domain (e.g. access) logic for reads
    private                  - this is private state, it can only be initialized directly from the repository
    required:"error message" - if not present in the constructor, an error with the provided message will be returned

  Expected Folder Structure:
    ./domain
    ├── livecall
    │   ├── livecall.go
    │   └── livecall_gen.go
    ├── party
    │   ├── party.go
    │   └── party_gen.go
    └── ...`,
	Example: `  Command:
    //go:generate go run github.com/xoe-labs/go-generators/ddd-domain-gen --type YOURTYPE
    ddd-domain-gen -t YOURTYPE

  Code:
    type Account struct {
        uuid    *string ` + "`" + `gen-getter,required:"field uuid is missing"` + "`" + `
        holder  *string ` + "`" + `gen-getter` + "`" + `
        balance *int64  ` + "`" + `private` + "`" + ` // reading the balance abides by domain logic
    }

    Required fields must be pointers for validation to work. So just use pointers everywhere.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generate.Main(sourceType)
	},
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ddd-domain-gen.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&sourceType, "type", "t", "", "The source type for which to generate the code")
	rootCmd.MarkFlagRequired("type")
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

		// Search config in home directory with name ".ddd-domain-gen" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ddd-domain-gen")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
