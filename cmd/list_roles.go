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
)

// listRolesCmd represents the list command
var listRolesCmd = &cobra.Command{
	Use:   "list",
	Short: "list roles defined in tenant",
	//	Long: `A longer description that spans multiple lines and likely contains examples
	//and usage of using your command. For example:
	//
	//Cobra is a CLI library for Go that empowers applications.
	//This application is a tool to generate the needed files
	//to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		listRoles()
	},
}

func init() {
	rolesCmd.AddCommand(listRolesCmd)
}

func listRoles() {
	rl, err := m.Role.List()
	if err != nil {
		exitWithMessage(err, 1)
	}

	ttx := table.NewWriter()
	ttx.AppendHeader(table.Row{"ID", "Name", "Description"})
	ttx.SetAutoIndex(true)
	for _, r := range rl.Roles {
		ttx.AppendRow(table.Row{*r.ID, *r.Name, *r.Description})
	}
	// Render output
	switch outFormat {
	case "csv":
		fmt.Println(ttx.RenderCSV())
	default:
		fmt.Println(ttx.Render())
	}
}
