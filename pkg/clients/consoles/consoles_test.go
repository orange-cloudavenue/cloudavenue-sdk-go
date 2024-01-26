package consoles

import (
	"testing"
)

func TestConsoles(t *testing.T) {
	tests := []struct {
		name    string
		console Console
		orgName string
		wantErr bool
	}{
		{
			name:    "should return the correct site name for Console1",
			console: Console1,
			orgName: "cav01ev01ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console2",
			console: Console2,
			orgName: "cav01iv02ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console4",
			console: Console4,
			orgName: "cav02ev04ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console5",
			console: Console5,
			orgName: "cav02iv05ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console7",
			console: Console7,
			orgName: "cav01iv07ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console8",
			console: Console8,
			orgName: "cav01iv08ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console9",
			console: Console9,
			orgName: "cav00vv09ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console9",
			console: Console9,
			orgName: "cav01vv09ocb0001234",
			wantErr: false,
		},
		{
			name:    "should return the correct site name for Console9",
			console: Console9,
			orgName: "cav02vv09ocb0001234",
			wantErr: false,
		},

		{
			name:    "should return an error if the organization is empty",
			console: "",
			orgName: "",
			wantErr: true,
		},
		{
			name:    "should return an error if the organization is invalid",
			console: "",
			orgName: "cav10ev01ocb0001234",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := FingByOrganizationName(tt.orgName)
			if (err != nil) && !tt.wantErr {
				t.Errorf("FingByOrganizationName(%s) error = %v, wantErr %v", tt.orgName, err, tt.wantErr)
				return
			}

			if !tt.wantErr && !CheckOrganizationName(tt.orgName) {
				t.Errorf("CheckOrganizationName(%s) error = %v, wantErr %v", tt.orgName, err, tt.wantErr)
				return
			}

			if c.GetSiteID() != tt.console {
				t.Errorf("FingByOrganizationName(%s) = %v, want %v", tt.orgName, c.GetSiteID(), tt.console)
				return
			}
		})
	}
}
