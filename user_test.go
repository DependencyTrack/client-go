package dtrack

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateManagedUser(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionAccessManagement,
		},
	})

	user, err := client.User.CreateManaged(context.Background(), ManagedUser{
		Username:        "test-managed",
		Fullname:        "test-managed-full-name",
		Email:           "test-managed-email@localhost",
		NewPassword:     "test-managed-password",
		ConfirmPassword: "test-managed-password",
	})
	require.NoError(t, err)

	require.Equal(t, user.Username, "test-managed")
	require.NotNil(t, user.LastPasswordChange)
	require.Equal(t, user.Fullname, "test-managed-full-name")
	require.Equal(t, user.Email, "test-managed-email@localhost")
	require.Equal(t, user.Teams, []Team{})
	require.Equal(t, user.Permissions, []Permission{})
	require.Equal(t, user.Suspended, false)
	require.Equal(t, user.ForcePasswordChange, false)
	require.Equal(t, user.NonExpiryPassword, false)
}

func TestUpdateManagedUser(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionAccessManagement,
		},
	})

	user, err := client.User.CreateManaged(context.Background(), ManagedUser{
		Username:        "test-managed",
		Fullname:        "test-managed-full-name",
		Email:           "test-managed-email@localhost",
		NewPassword:     "test-managed-password",
		ConfirmPassword: "test-managed-password",
	})
	require.NoError(t, err)
	require.Equal(t, user.Username, "test-managed")
	require.Equal(t, user.Fullname, "test-managed-full-name")
	require.Equal(t, user.Email, "test-managed-email@localhost")

	updated, err := client.User.UpdateManaged(context.Background(), ManagedUser{
		Username: user.Username,
		Fullname: "test-managed-full-name-updated",
		Email:    user.Email,
	})

	require.NoError(t, err)
	require.Equal(t, updated.Username, user.Username)
	require.Equal(t, updated.Fullname, "test-managed-full-name-updated")
	require.Equal(t, updated.Email, user.Email)
}

func TestGetAllManagedUsers(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionAccessManagement,
		},
	})

	user, err := client.User.CreateManaged(context.Background(), ManagedUser{
		Username:        "test-managed",
		Fullname:        "test-managed-full-name",
		Email:           "test-managed-email@localhost",
		NewPassword:     "test-managed-password",
		ConfirmPassword: "test-managed-password",
	})
	require.NoError(t, err)
	require.Equal(t, user.Username, "test-managed")
	require.Equal(t, user.Fullname, "test-managed-full-name")
	require.Equal(t, user.Email, "test-managed-email@localhost")

	users, err := FetchAll(func(po PageOptions) (Page[ManagedUser], error) {
		return client.User.GetAllManaged(context.Background(), po)
	})

	require.NoError(t, err)
	require.Equal(t, len(users), 2)
	require.Equal(t, users[0].Username, "admin")
	require.Equal(t, len(users[0].Teams), 1)
	require.Equal(t, users[0].Teams[0].Name, "Administrators")
	require.Equal(t, len(users[0].Permissions), 14)

	require.Equal(t, users[1].Username, "test-managed")
}

func TestDeleteManagedUser(t *testing.T) {
	// Create User
	// Validate creation of user
	// Delete user
	// Validate deletion of user
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionAccessManagement,
		},
	})

	user, err := client.User.CreateManaged(context.Background(), ManagedUser{
		Username:        "test-managed",
		Fullname:        "test-managed-full-name",
		Email:           "test-managed-email@localhost",
		NewPassword:     "test-managed-password",
		ConfirmPassword: "test-managed-password",
	})
	require.NoError(t, err)
	require.Equal(t, user.Username, "test-managed")
	require.Equal(t, user.Fullname, "test-managed-full-name")
	require.Equal(t, user.Email, "test-managed-email@localhost")

	users, err := FetchAll(func(po PageOptions) (Page[ManagedUser], error) {
		return client.User.GetAllManaged(context.Background(), po)
	})

	require.NoError(t, err)
	require.Equal(t, len(users), 2)

	err = client.User.DeleteManaged(context.Background(), ManagedUser{
		Username: user.Username,
	})

	require.NoError(t, err)

	users, err = FetchAll(func(po PageOptions) (Page[ManagedUser], error) {
		return client.User.GetAllManaged(context.Background(), po)
	})

	require.NoError(t, err)
	require.Equal(t, len(users), 1)
}
