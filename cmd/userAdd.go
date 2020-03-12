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
	"gopkg.in/auth0.v3/management"
	"os"
	"strings"
)

var password, investorID, supplierID string
var roles []string

// userAddCmd represents the userAdd command
var userAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add new user",
	Long:  `Adds a new user`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email := strings.ToLower(args[0])
		addUser(&email, &password, &investorID, &supplierID)
	},
}

func init() {
	usersCmd.AddCommand(userAddCmd)
	userAddCmd.Flags().StringVarP(&password, "password", "p", "", "user password")
	userAddCmd.Flags().StringVarP(&investorID, "inverstorid", "i", "", "ID of investor")
	userAddCmd.Flags().StringVarP(&supplierID, "supplierid", "s", "", "ID of supplier")
	userAddCmd.Flags().StringSliceVarP(&roles, "roles", "r", []string{}, "list of roles")
}

func addUser(email, password, investorID, supplierID *string) {
	appMetadata := map[string]interface{}{"investorId": investorID, "supplierId": supplierID}
	userMetadata := map[string]interface{}{}
	conn := "Username-Password-Authentication"

	var u = &management.User{
		Name:         email,
		Email:        email,
		Connection:   &conn,
		Password:     password,
		UserMetadata: userMetadata,
		AppMetadata:  appMetadata,
	}

	if err := m.User.Create(u); err != nil {
		if err, ok := err.(management.Error); ok {
			// https://auth0.com/docs/api/management/v2#!/Users/post_users
			switch err.Status() {
			case 409:
				if ul, err := m.User.ListByEmail(*email); err != nil {
					fmt.Println(err)
				} else {
					for _, u := range ul {
						fmt.Printf("User with email %s already exists with ID %s\n", *u.Email, *u.ID)
					}
				}
			case 400:
				if strings.Contains(err.Error(), "PasswordStrengthError: Password is too weak") {
					fmt.Printf("Cannot create user %s due to %s\n", *u.Name, err.Error())
				} else {
					fmt.Println(err.Error())
				}
			default:
				fmt.Println(err)
			}
		}
		os.Exit(1)
	} else {
		fmt.Printf("user %s has been created with ID: %s\n", *u.Name, *u.ID)
	}

	if err := addRolesToUser(u, roles...); err != nil {
		fmt.Println(err)
	}
}

func addRolesToUser(u *management.User, roles ...string) error {
	availableRoles, err := getRoleMap()
	if err != nil {
		return err
	}

	var rolesToAdd []*management.Role
	for _, role := range roles {
		if _, ok := availableRoles[role]; ok {
			rolesToAdd = append(rolesToAdd, availableRoles[role])
		} else {
			fmt.Printf("Role '%s' doesn't exist\n", role)
		}

	}
	if err := m.User.AssignRoles(*u.ID, rolesToAdd...); err != nil {
		return err
	}
	return nil
}

func getRoleMap() (map[string]*management.Role, error) {
	var roleMap = make(map[string]*management.Role)
	r, err := m.Role.List()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, role := range r.Roles {
		roleMap[*role.Name] = role
	}

	return roleMap, nil
}
