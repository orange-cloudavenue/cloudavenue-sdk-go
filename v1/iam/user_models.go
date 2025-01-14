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

//go:generate mockgen -source=user_models.go -destination=mock/zz_generated_user_client.go

type (
	UserClient struct {
		govcdAdminOrg clientGoVCDAdminOrg
		govcdUser     *govcd.OrgUser

		// Data
		User *User
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
		Password string `validate:"required,gte=6"`
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
	for i := 0; i < x.NumField(); i++ {
		field := x.Type().Field(i)
		value := x.Field(i).Interface()
		if field.Name == "User" {
			// If the field is a User struct, iterate over its fields
			for j := 0; j < reflect.ValueOf(value).NumField(); j++ {
				userField := reflect.ValueOf(value).Type().Field(j)
				userValue := reflect.ValueOf(value).Field(j).Interface()

				switch userField.Name {
				case "Name":
					u.Name = userValue.(string)
				case "RoleName":
					u.Role = roleReference
				case "Description":
					u.Description = userValue.(string)
				case "FullName":
					u.FullName = userValue.(string)
				case "Email":
					u.EmailAddress = userValue.(string)
				case "Telephone":
					u.Telephone = userValue.(string)
				case "Enabled":
					u.IsEnabled = userValue.(bool)
				case "DeployedVMQuota":
					u.DeployedVmQuota = userValue.(int)
				case "StoredVMQuota":
					u.StoredVmQuota = userValue.(int)
				case "ID":
					u.ID = userValue.(string)
				}
			}
		}

		if field.Name == "Password" {
			u.Password = value.(string)
		}
	}

	// Switch name of the user struct to set specific fields
	// No default case because only LocalUser and SAMLUser are extra fields
	switch v := user.(type) {
	case LocalUser:
		u.Password = v.Password
		u.ProviderType = govcd.OrgUserProviderIntegrated
	case SAMLUser:
		u.ProviderType = govcd.OrgUserProviderSAML
		u.IsExternal = true
	}

	return u
}

func toSDKTypeUser[toUser any](user *govcdtypes.User) *toUser {
	u := new(toUser)

	// toSDKTypeUser converts a go-vcd user to a user by reflecting.
	x := reflect.ValueOf(u).Elem()
	for i := 0; i < x.NumField(); i++ {
		field := x.Type().Field(i)

		switch field.Name {
		case "Name":
			x.Field(i).SetString(user.Name)
		case "RoleName":
			x.Field(i).SetString(user.Role.Name)
		case "Description":
			x.Field(i).SetString(user.Description)
		case "FullName":
			x.Field(i).SetString(user.FullName)
		case "Email":
			x.Field(i).SetString(user.EmailAddress)
		case "Telephone":
			x.Field(i).SetString(user.Telephone)
		case "Enabled":
			x.Field(i).SetBool(user.IsEnabled)
		case "DeployedVMQuota":
			x.Field(i).SetInt(int64(user.DeployedVmQuota))
		case "StoredVMQuota":
			x.Field(i).SetInt(int64(user.StoredVmQuota))
		case "ID":
			x.Field(i).SetString(user.ID)
		}
	}

	// Set specific fields based on the ProviderType
	switch user.ProviderType {
	case govcd.OrgUserProviderIntegrated:
		x.FieldByName("Type").SetString(string(UserTypeLocal))
	case govcd.OrgUserProviderSAML:
		x.FieldByName("Type").SetString(string(UserTypeSAML))
	}

	return u
}
