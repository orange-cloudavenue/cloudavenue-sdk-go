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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/helpers/validators"
)

// CreateLocalUser creates a new local user in the system.
//
// Parameters:
//
//	user - The LocalUser struct containing the details of the user to be created.
//
// Returns:
//
//	*User - A pointer to the created User struct.
//	error - An error if the creation or validation fails, otherwise nil.
func (c *Client) CreateLocalUser(user LocalUser) (*UserClient, error) {
	return c.createGenericUser(user)
}

// CreateSAMLUser creates a new SAML user in the system.
// It takes a SAMLUser object as input and returns a pointer to the created User object and an error, if any.
//
// Parameters:
//   - user: SAMLUser object containing the details of the user to be created.
//
// Returns:
//   - *User: Pointer to the created User object.
//   - error: Error, if any occurred during the creation process.
func (c *Client) CreateSAMLUser(user SAMLUser) (*UserClient, error) {
	return c.createGenericUser(user)
}

func (c *Client) createGenericUser(user userInterface) (*UserClient, error) {
	// Refresh the client if needed
	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	// Validate struct with tags defined in the struct itself
	if err := validators.New().Struct(user); err != nil {
		return nil, err
	}

	// Get Role HREF
	roleRef, err := c.clientGoVCDAdminOrg.GetRoleReference(user.GetRoleName())
	if err != nil {
		return nil, err
	}

	// Create the user in the system
	userCreated, err := c.clientGoVCDAdminOrg.CreateUser(toGoVCDTypeUser(user, roleRef))
	if err != nil {
		return nil, err
	}

	return &UserClient{
		govcdAdminOrg: c.clientGoVCDAdminOrg,
		govcdUser:     userCreated,
		User:          toSDKTypeUser(userCreated.User),
	}, nil
}

// GetUser retrieves a user by their name or ID.
// GetUser permit to retrieve a local or SAML user by their name or ID.
//
// Parameters:
//   - nameOrID: The name or ID of the user to retrieve.
//
// Returns:
//   - A pointer to the User object if found, or an error if any issues occur during the process.
func (c *Client) GetUser(nameOrID string) (*UserClient, error) {
	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	user, err := c.clientGoVCDAdminOrg.GetUserByNameOrId(nameOrID, true)
	if err != nil {
		return nil, err
	}

	return &UserClient{
		govcdAdminOrg: c.clientGoVCDAdminOrg,
		govcdUser:     user,
		User:          toSDKTypeUser(user.User),
	}, nil
}

// Update updates the user information in the system.
// Returns an error if the update operation fails.
func (u *UserClient) Update() error {
	if err := validators.New().Struct(u.User); err != nil {
		return err
	}

	old := *u.govcdUser.User

	// Get Role HREF
	roleRef, err := u.govcdAdminOrg.GetRoleReference(u.User.RoleName)
	if err != nil {
		return err
	}

	u.govcdUser.User = toGoVCDTypeUser(u.User, roleRef)
	u.govcdUser.User.Href = old.Href
	u.govcdUser.User.ProviderType = old.ProviderType
	u.govcdUser.User.IsExternal = old.IsExternal
	u.govcdUser.User.ID = old.ID

	// Update the user
	return u.govcdUser.Update()
}

// Delete deletes a user from the system.
func (u *UserClient) Delete(takeOwnership bool) error {
	return u.govcdUser.Delete(takeOwnership)
}

// Enable enables a user if it was disabled. Fails otherwise.
func (u *UserClient) Enable() error {
	return u.govcdUser.Enable()
}

// Disable disables a user if it was enabled. Fails otherwise.
func (u *UserClient) Disable() error {
	return u.govcdUser.Disable()
}

// Unlock unlocks a user if it was locked. Fails otherwise.
func (u *UserClient) Unlock() error {
	return u.govcdUser.Unlock()
}

// ChangePassword changes the password of a user.
func (u *UserClient) ChangePassword(password string) error {
	return u.govcdUser.ChangePassword(password)
}
