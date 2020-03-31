// Copyright © 2020 author from config
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
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"reflect"
)

// userGetCmd represents the userGet command
var userGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "get user by id",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		getUser(args[0])
	},
}

func init() {
	usersCmd.AddCommand(userGetCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getUser(id string) {

	u, err := m.User.Read(id)
	if err != nil {
		exitWithMessage(err, 0)
	}

	ttx := table.NewWriter()
	ttx.AppendHeader(table.Row{"ID", "Name", "Email", "Nickname", "Last login"})
	ttx.SetAutoIndex(true)
	ttx.AppendRow(table.Row{
		getVal(u.ID),
		getVal(u.Name),
		getVal(u.Email),
		getVal(u.Nickname),
		getVal(u.LastLogin),
		})

	// Render output
	switch outFormat {
	case "csv":
		fmt.Println(ttx.RenderCSV())
	default:
		fmt.Println(ttx.Render())
	}
}


func getVal(v interface{}) interface{} {
	if !reflect.ValueOf(v).IsNil() {
		return reflect.ValueOf(v).Elem()
	}
	return "N/A"
}
