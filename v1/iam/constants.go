package iam

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

const (
	// UserTypeLocal is the type of the user.
	UserTypeLocal UserType = govcd.OrgUserProviderIntegrated
	// UserTypeSAML is the type of the user.
	UserTypeSAML UserType = govcd.OrgUserProviderSAML
)
