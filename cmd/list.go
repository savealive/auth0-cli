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
	"github.com/savealive/auth0-cli/manager"
	"github.com/spf13/cobra"
	"gopkg.in/auth0.v3/management"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		m, err := manager.New()
		if err != nil {
			panic(err)
		}
		List(m)
	},
}

var filterRole string

type auth0User struct {
	user  *management.User
	roles []string
}

type rolesMap map[management.Role][]*management.User

func init() {

	usersCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listCmd.Flags().StringVarP(&filterRole, "role", "r", "", "Role to filter")
}

func GetRoles(m *management.Management) (map[string]string, error) {
	roles := make(map[string]string)
	rl, err := m.Role.List()
	if err != nil {
		return nil, err
	}
	for _, r := range rl.Roles {
		roles[*r.Name] = *r.ID
	}
	return roles, nil
}

func getRolesWithUsers(m *management.Management) (rolesMap, error) {
	r := make(map[management.Role][]*management.User)
	roles, err := m.Role.List()
	if err != nil {
		return nil, err
	}
	for _, role := range roles.Roles {
		users, err := m.Role.Users(*role.ID)
		if err != nil {
			return nil, err
		}
		r[*role] = users
	}
	return r, nil
}

func List(m *management.Management) {
	printUsers(m, filterRole)
}

func fetchAllUsers(m *management.Management) ([]*management.User, error) {

	l, err := m.User.List()
	if err != nil {
		return nil, err
	}
	users := make([]*management.User, 0, l.Total)
	users = append(users, l.Users...)
	for page := 1; l.HasNext(); page++ {
		l, err = m.User.List(management.Page(page))
		if err != nil {
			return nil, err
		}
		users = append(users, l.Users...)
	}
	if err != nil {
		return nil, err
	}
	return users, nil
}

func appendRolesToUser(users []*management.User, r rolesMap) []auth0User {
	res := make([]auth0User, 0, len(users))
	for _, u := range users {
		var user = auth0User{
			user: u,
		}
		for k, v := range r {
			if stringInUsersSlice(v, u) {
				user.roles = append(user.roles, *k.Name)
			}
		}
		res = append(res, user)
	}
	return res
}

func stringInUsersSlice(s []*management.User, u *management.User) bool {
	for _, user := range s {
		if *user.ID == *u.ID {
			return true
		}
	}
	return false
}

func printUsers(m *management.Management, role string) {
	ttx := table.NewWriter()
	ttx.AppendHeader(table.Row{"ID", "User", "Roles"})
	ttx.SetAutoIndex(true)
	var users []*management.User
	var err error
	if role != "" {
		users, err = fetchUsersByRole(m, role)

	} else {
		users, err = fetchAllUsers(m)
	}
	if err != nil {
		panic(err)
	}

	rolesUsers, err := getRolesWithUsers(m)
	if err != nil {
		panic(err)
	}

	usersWithRoles := appendRolesToUser(users, rolesUsers)

	for _, u := range usersWithRoles {
		ttx.AppendRow(table.Row{*u.user.ID, *u.user.Name, u.roles})
	}
	fmt.Println(ttx.Render())
}

func fetchUsersByRole(m *management.Management, r string) ([]*management.User, error) {
	roles, err := GetRoles(m)
	if err != nil {
		return nil, err
	}
	return m.Role.Users(roles[r])
}
