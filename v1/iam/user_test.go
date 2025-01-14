package iam

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	mock "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/iam/mock"
)

func TestClient_CreateLocalUser(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for admin org.
	clientAdminOrg := mock.NewMockclientGoVCDAdminOrg(ctrl)

	// Mock client for cloudavenue.
	clientCAV := mock.NewMockclientCloudavenue(ctrl)

	c := &Client{
		clientGoVCDAdminOrg: clientAdminOrg,
		clientCloudavenue:   clientCAV,
	}

	testCases := []struct {
		name              string
		mockFunc          []func()
		userValue         LocalUser
		expectedUserValue *LocalUser
		expectedErr       bool
		err               error
	}{
		{
			name: "success",
			mockFunc: []func(){
				func() {
					clientAdminOrg.EXPECT().CreateUser(gomock.Any()).Return(&govcd.OrgUser{
						User: &govcdtypes.User{
							Name: "test",
							Role: &govcdtypes.Reference{
								Name: "test",
							},
							ID: urn.VcloudPrefix + urn.User.String() + uuid.NewString(),
						},
					}, nil)
				},
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
				func() {
					clientAdminOrg.EXPECT().GetRoleReference(gomock.Any()).Return(&govcdtypes.Reference{
						Name: "test",
						HREF: "https://foo.bar/api/admin/role/test",
					}, nil)
				},
			},
			userValue: LocalUser{
				User: User{
					Name:     "test",
					RoleName: "test",
				},

				Password: "test",
			},
			expectedUserValue: &LocalUser{
				User: User{
					Name:     "test",
					RoleName: "test",
					ID:       urn.VcloudPrefix + urn.User.String(),
				},
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "error-validation",
			mockFunc: []func(){
				// func CreateUser is not called because the validation fails.
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
			},
			userValue: LocalUser{
				User: User{
					Name:     "test",
					RoleName: "test",
				},
				Password: "",
			},
			expectedErr: true,
			err:         errors.New("'LocalUser.Password' Error:Field validation for 'Password' failed on the 'required' tag"),
		},
		{
			name: "error-refresh",
			mockFunc: []func(){
				// func CreateUser is not called because the refresh fails.
				func() {
					clientCAV.EXPECT().Refresh().Return(errors.New("error"))
				},
			},
			userValue:   LocalUser{},
			expectedErr: true,
		},
		{
			name: "error-create-user",
			mockFunc: []func(){
				func() {
					clientAdminOrg.EXPECT().CreateUser(gomock.Any()).Return(nil, errors.New("error"))
				},
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
				func() {
					clientAdminOrg.EXPECT().GetRoleReference(gomock.Any()).Return(&govcdtypes.Reference{
						Name: "test",
						HREF: "https://foo.bar/api/admin/role/test",
					}, nil)
				},
			},
			userValue: LocalUser{
				User: User{
					Name:     "test",
					RoleName: "test",
				},
				Password: "1234",
			},
			expectedErr: true,
		},
		{
			name: "error-get-role",
			mockFunc: []func(){
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
				func() {
					clientAdminOrg.EXPECT().GetRoleReference(gomock.Any()).Return(nil, errors.New("error"))
				},
			},
			userValue: LocalUser{
				User: User{
					Name:     "test",
					RoleName: "test",
				},
				Password: "1234",
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, f := range tc.mockFunc {
				f()
			}

			user, err := c.CreateLocalUser(tc.userValue)
			if tc.expectedErr && err == nil {
				t.Errorf("expected error %v, got %v", tc.err, err)
			}

			if !tc.expectedErr && err != nil {
				t.Errorf("no expected error %v, got %v", tc.err, err)
			}

			if (tc.expectedErr && tc.err != nil) && err != nil {
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Errorf("expected error %v, got %v", tc.err, err)
				}
			}

			if tc.expectedUserValue != nil {
				if !strings.Contains(user.User.ID, tc.expectedUserValue.User.ID) {
					t.Errorf("expected user ID prefix %v, got %v", urn.VcloudPrefix+urn.User.String(), user.User.ID)
				}
			}
		})
	}
}

func TestClient_CreateSAMLUser(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for admin org.
	clientAdminOrg := mock.NewMockclientGoVCDAdminOrg(ctrl)

	// Mock client for cloudavenue.
	clientCAV := mock.NewMockclientCloudavenue(ctrl)

	c := &Client{
		clientGoVCDAdminOrg: clientAdminOrg,
		clientCloudavenue:   clientCAV,
	}

	testCases := []struct {
		name              string
		mockFunc          []func()
		userValue         SAMLUser
		expectedUserValue *SAMLUser
		expectedErr       bool
		err               error
	}{
		{
			name: "success",
			mockFunc: []func(){
				func() {
					clientAdminOrg.EXPECT().CreateUser(gomock.Any()).Return(&govcd.OrgUser{
						User: &govcdtypes.User{
							Name: "test",
							Role: &govcdtypes.Reference{
								Name: "test",
							},
							ID: urn.VcloudPrefix + urn.User.String() + uuid.NewString(),
						},
					}, nil).AnyTimes()
				},
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
				func() {
					clientAdminOrg.EXPECT().GetRoleReference(gomock.Any()).Return(&govcdtypes.Reference{
						Name: "test",
						HREF: "https://foo.bar/api/admin/role/test",
					}, nil)
				},
			},
			userValue: SAMLUser{
				User: User{
					Name:     "test",
					RoleName: "test",
				},
			},
			expectedUserValue: &SAMLUser{
				User: User{
					Name:     "test",
					RoleName: "test",
					ID:       urn.VcloudPrefix + urn.User.String(),
				},
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "error-validation",
			mockFunc: []func(){
				// func CreateUser is not called because the validation fails.
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
			},
			userValue: SAMLUser{
				User: User{
					Name:     "test",
					RoleName: "",
				},
			},
			expectedErr: true,
			err:         errors.New("Error:Field validation for 'RoleName' failed on the 'required' tag"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, f := range tc.mockFunc {
				f()
			}

			user, err := c.CreateSAMLUser(tc.userValue)
			if tc.expectedErr && err == nil {
				t.Errorf("expected error %v, got %v", tc.err, err)
			}

			if !tc.expectedErr && err != nil {
				t.Errorf("no expected error %v, got %v", tc.err, err)
			}

			if (tc.expectedErr && tc.err != nil) && err != nil {
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Errorf("expected error %v, got %v", tc.err, err)
				}
			}

			if tc.expectedUserValue != nil {
				if !strings.Contains(user.User.ID, tc.expectedUserValue.User.ID) {
					t.Errorf("expected user ID prefix %v, got %v", urn.VcloudPrefix+urn.User.String(), user.User.ID)
				}
			}
		})
	}
}

func TestClient_GetUser(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for admin org.
	clientAdminOrg := mock.NewMockclientGoVCDAdminOrg(ctrl)

	// Mock client for cloudavenue.
	clientCAV := mock.NewMockclientCloudavenue(ctrl)

	c := &Client{
		clientGoVCDAdminOrg: clientAdminOrg,
		clientCloudavenue:   clientCAV,
	}

	testCases := []struct {
		name              string
		mockFunc          []func()
		nameOrID          string
		expectedUserValue *User
		expectedErr       bool
		err               error
	}{
		{
			name: "success-local",
			mockFunc: []func(){
				func() {
					clientAdminOrg.EXPECT().GetUserByNameOrId(gomock.Any(), true).Return(&govcd.OrgUser{
						User: &govcdtypes.User{
							Name: "test",
							Role: &govcdtypes.Reference{
								Name: "test",
							},
							ProviderType: govcd.OrgUserProviderIntegrated,
							ID:           urn.VcloudPrefix + urn.User.String() + uuid.NewString(),
						},
					}, nil)
				},
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
			},
			nameOrID: "test",
			expectedUserValue: &User{
				Name:     "test",
				RoleName: "test",
				ID:       urn.VcloudPrefix + urn.User.String(),
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "success-saml",
			mockFunc: []func(){
				func() {
					clientAdminOrg.EXPECT().GetUserByNameOrId(gomock.Any(), true).Return(&govcd.OrgUser{
						User: &govcdtypes.User{
							Name: "test",
							Role: &govcdtypes.Reference{
								Name: "test",
							},
							ProviderType: govcd.OrgUserProviderSAML,
							ID:           urn.VcloudPrefix + urn.User.String() + uuid.NewString(),
						},
					}, nil)
				},
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
			},
			nameOrID: "test",
			expectedUserValue: &User{
				Name:     "test",
				RoleName: "test",
				ID:       urn.VcloudPrefix + urn.User.String(),
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "error-refresh",
			mockFunc: []func(){
				// func GetUserByNameOrId is not called because the refresh fails.
				func() {
					clientCAV.EXPECT().Refresh().Return(errors.New("error"))
				},
			},
			nameOrID:    "test",
			expectedErr: true,
		},
		{
			name: "error-get-user",
			mockFunc: []func(){
				func() {
					clientAdminOrg.EXPECT().GetUserByNameOrId(gomock.Any(), true).Return(nil, errors.New("error"))
				},
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
			},
			nameOrID:    "test",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, f := range tc.mockFunc {
				f()
			}

			user, err := c.GetUser(tc.nameOrID)
			if tc.expectedErr && err == nil {
				t.Errorf("expected error %v, got %v", tc.err, err)
			}

			if !tc.expectedErr && err != nil {
				t.Errorf("no expected error %v, got %v", tc.err, err)
			}

			if (tc.expectedErr && tc.err != nil) && err != nil {
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Errorf("expected error %v, got %v", tc.err, err)
				}
			}

			if tc.expectedUserValue != nil {
				if !strings.Contains(user.User.ID, tc.expectedUserValue.ID) {
					t.Errorf("expected user ID prefix %v, got %v", urn.VcloudPrefix+urn.User.String(), user.User.ID)
				}
			}
		})
	}
}
