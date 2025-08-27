/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iam

import (
	"reflect"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

var (
	_ userInterface = LocalUser{}
	_ userInterface = SAMLUser{}
)

//go:generate mockgen -source=user_types.go -destination=mock/zz_generated_user_client.go

type (
	UserClient struct {
		govcdAdminOrg clientGoVCDAdminOrg
		govcdUser     *govcd.OrgUser

		// Data
		User User
	}

	userInterface interface {
		GetRoleName() string
	}

	User struct {
		// REQUIRED: The name of the user.
		Name     string `validate:"required,disallow_upper,disallow_space"`
		RoleName string `validate:"required"`

		// OPTIONAL
		Description     string `validate:"omitempty"`
		FullName        string `validate:"omitempty"`
		Email           string `validate:"omitempty,email"`
		Telephone       string `validate:"omitempty"`
		Enabled         bool
		DeployedVMQuota int // 0 means unlimited
		StoredVMQuota   int // 0 means unlimited

		// READ-ONLY
		ID   string
		Type UserType
	}

	// UserType is a type of user.
	UserType string

	LocalUser struct {
		User `validate:"required"`

		// REQUIRED: The password of the user.
		Password string `validate:"required,min=6"`
	}

	SAMLUser struct {
		User `validate:"required"`
	}
)

// GetRoleName returns the role name of the user.
func (u User) GetRoleName() string {
	return u.RoleName
}

// toGoVCDTypeUser converts a user to a go-vcd user by reflecting on the fields of the input user struct.
// It supports conversion from LocalUser and SAMLUser types to govcdtypes.User.
//
// Parameters:
//   - user: an input user of any type (expected to be either LocalUser or SAMLUser).
//
// Returns:
//   - *govcdtypes.User: a pointer to the converted govcdtypes.User struct.
//   - error: an error if the input user type is not recognized.
//
// The function uses reflection to iterate over the fields of the input user struct and map them to the corresponding
// fields in the govcdtypes.User struct. It handles specific fields such as Name, RoleName, Description, FullName,
// Email, Telephone, Enabled, DeployedVMQuota, StoredVMQuota, and ID. Additionally, it sets the Password and ProviderType
// fields based on the specific user type (LocalUser or SAMLUser).
func toGoVCDTypeUser(user any, roleReference *govcdtypes.Reference) *govcdtypes.User {
	// toGoVCDTypeUser converts a user to a go-vcd user by reflecting.
	x := reflect.ValueOf(user)
	if x.Kind() == reflect.Ptr {
		x = x.Elem()
	}

	u := &govcdtypes.User{}

	sw := func(fieldName string, fieldValue any) {
		switch fieldName {
		case "Name":
			u.Name = fieldValue.(string)
		case "RoleName":
			u.Role = roleReference
		case "Description":
			u.Description = fieldValue.(string)
		case "FullName":
			u.FullName = fieldValue.(string)
		case "Email":
			u.EmailAddress = fieldValue.(string)
		case "Telephone":
			u.Telephone = fieldValue.(string)
		case "Enabled":
			u.IsEnabled = fieldValue.(bool)
		case "DeployedVMQuota":
			u.DeployedVmQuota = fieldValue.(int)
		case "StoredVMQuota":
			u.StoredVmQuota = fieldValue.(int)
		case "ID":
			u.ID = fieldValue.(string)
		case "Password":
			u.Password = fieldValue.(string)
		}
	}

	for i := 0; i < x.NumField(); i++ {
		field := x.Type().Field(i)
		value := x.Field(i).Interface()

		switch field.Name {
		case "User":
			userValue := reflect.ValueOf(value)
			for j := 0; j < userValue.NumField(); j++ {
				userField := userValue.Type().Field(j)
				userFieldValue := userValue.Field(j).Interface()

				sw(userField.Name, userFieldValue)
			}
		default:
			sw(field.Name, value)
		}
	}

	// Switch name of the user struct to set specific fields
	// No default case because only LocalUser and SAMLUser are extra fields
	switch v := user.(type) {
	case *LocalUser:
		u.Password = v.Password
		u.ProviderType = govcd.OrgUserProviderIntegrated
	case *SAMLUser:
		u.ProviderType = govcd.OrgUserProviderSAML
		u.IsExternal = true
	}

	u.Xmlns = govcdtypes.XMLNamespaceVCloud
	u.Type = govcdtypes.MimeAdminUser

	return u
}

func toSDKTypeUser(user *govcdtypes.User) User {
	u := User{
		Name:            user.Name,
		RoleName:        user.Role.Name,
		Description:     user.Description,
		FullName:        user.FullName,
		Email:           user.EmailAddress,
		Telephone:       user.Telephone,
		Enabled:         user.IsEnabled,
		DeployedVMQuota: user.DeployedVmQuota,
		StoredVMQuota:   user.StoredVmQuota,
		ID:              user.ID,
	}

	// Set specific fields based on the ProviderType
	switch user.ProviderType {
	case govcd.OrgUserProviderIntegrated:
		u.Type = UserTypeLocal
	case govcd.OrgUserProviderSAML:
		u.Type = UserTypeSAML
	}

	return u
}
