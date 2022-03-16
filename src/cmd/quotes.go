/*
Copyright © 2022 Joaquin Genova <joaquingenovag8@gmail.com>

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

	"github.com/spf13/cobra"
)

// quotesCmd represents the quotes command
var quotesCmd = &cobra.Command{
	Use:     "quotes",
	Aliases: []string{"qut"},
	Short:   "Return the path between quotes",

	Run: func(cmd *cobra.Command, args []string) {

		if path, err := CheckIndexOrAbbvOrDir(args[0]); err == nil {
			fmt.Println("\"" + path + "\"")
			os.Exit(0)
		} else {
			cobra.CheckErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(quotesCmd)
}
