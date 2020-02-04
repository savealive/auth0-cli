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
	"github.com/spf13/cobra"
)

// userDeleteCmd represents the userDelete command
var userDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete user by id or email",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if err := deleteUser(id); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("user id: %s deleted\n", id)
		}
	},
}

func init() {
	usersCmd.AddCommand(userDeleteCmd)
}

func userExists(id string) error {
	_, err := m.User.Read(id)
	return err
}

func deleteUser(id string) error {
	if err := userExists(id); err != nil {
		return fmt.Errorf("id: %s %v", id, err)
	}
	return m.User.Delete(id)
}
