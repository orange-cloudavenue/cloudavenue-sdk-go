package netbackup

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	clientnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/netbackup"
	commonnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/netbackup"
)

type getProtectionLevelResponse struct {
	Data ProtectionLevel `json:"data,omitempty"`
}

type getProtectionLevelByIDRequest struct {
	// Use one of the following to specify the protection level
	VAppID    *int
	VDCID     *int
	MachineID *int

	ProtectionLevelID *int
}

// GetProtectionLevel - Get a protection level by ID
// Use GetProtectionLevelsRequest to specify the VAppID, VDCID or MachineID and the ProtectionLevelID
func (p *ProtectionLevelClient) getProtectionLevelByID(req getProtectionLevelByIDRequest) (resp *ProtectionLevel, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return resp, err
	}

	if req.VAppID == nil && req.VDCID == nil && req.MachineID == nil {
		return resp, fmt.Errorf("you must specify a VAppID, VDCID or MachineID")
	}

	if req.ProtectionLevelID == nil {
		return resp, fmt.Errorf("you must specify a ProtectionLevelID")
	}

	var r *resty.Response

	hReq := c.R().
		SetResult(&getProtectionLevelResponse{}).
		SetError(&commonnetbackup.APIError{})

	switch {
	case req.VAppID != nil:
		r, err = hReq.
			SetPathParams(map[string]string{
				"VAppId":            fmt.Sprintf("%d", *req.VAppID),
				"ProtectionLevelId": fmt.Sprintf("%d", *req.ProtectionLevelID),
			}).
			Get("/v6/vcloud/vapps/{VAppId}/protection/levels/{ProtectionLevelId}")
	case req.VDCID != nil:
		r, err = hReq.
			SetPathParams(map[string]string{
				"VdcId":             fmt.Sprintf("%d", *req.VDCID),
				"ProtectionLevelId": fmt.Sprintf("%d", *req.ProtectionLevelID),
			}).
			Get("/v6/vcloud/vdcs/{VdcId}/protection/levels/{ProtectionLevelId}")
	case req.MachineID != nil:
		r, err = hReq.
			SetPathParams(map[string]string{
				"MachineId":         fmt.Sprintf("%d", *req.MachineID),
				"ProtectionLevelId": fmt.Sprintf("%d", *req.ProtectionLevelID),
			}).
			Get("/v6/machines/{MachineId}/protection/levels/{ProtectionLevelId}")
	}
	if err != nil {
		return resp, err
	}

	if r.IsError() {
		return resp, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return &r.Result().(*getProtectionLevelResponse).Data, nil
}

type getProtectionLevelByNameRequest struct {
	// Use one of the following to specify the protection level
	VAppID    *int `json:"VAppId,omitempty"`
	VDCID     *int `json:"VdcId,omitempty"`
	MachineID *int `json:"MachineId,omitempty"`

	ProtectionLevelName *string `json:"ProtectionLevelName"`
}

// GetProtectionLevelByName - Get a protection level by Name
// Use GetProtectionLevelsRequest to specify the VAppID, VDCID or MachineID and the ProtectionLevelName
func (p *ProtectionLevelClient) getProtectionLevelByName(req getProtectionLevelByNameRequest) (resp *ProtectionLevel, err error) {
	if req.VAppID == nil && req.VDCID == nil && req.MachineID == nil {
		return resp, fmt.Errorf("you must specify a VAppID, VDCID or MachineID")
	}

	if req.ProtectionLevelName == nil {
		return resp, fmt.Errorf("you must specify a ProtectionLevelName")
	}

	protectionLevels, err := p.ListProtectionLevels(listProtectionLevelsRequest{
		VAppID:    req.VAppID,
		VDCID:     req.VDCID,
		MachineID: req.MachineID,
	})
	if err != nil {
		return resp, err
	}

	for _, protectionLevel := range *protectionLevels {
		if protectionLevel.Name == *req.ProtectionLevelName {
			return &protectionLevel, nil
		}
	}

	return resp, fmt.Errorf("protection level %s not found", *req.ProtectionLevelName)
}
