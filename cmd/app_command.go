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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xoe-labs/ddd-gen/pkg/gen_app/command"
)

// appCommandCmd represents the app command
var appCommandCmd = &cobra.Command{
	Use:   "command",
	Short: "Generates idiomatic go code for commands (CQRS) within the application layer",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := command.NewConfig(
			viper.GetString("aggregate"),
			viper.GetString("policeable"),
			viper.GetString("identifiable"),
			viper.GetString("repository"),
			viper.GetString("policer"),
		)
		if err != nil {
			return err
		}
		return command.Gen(sourceType, *cfg)
	},
}

func init() {
	appCmd.AddCommand(appCommandCmd)
}
