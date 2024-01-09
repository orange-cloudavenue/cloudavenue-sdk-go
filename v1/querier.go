package v1

import (
	"fmt"
	"regexp"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
)

type (
	Query struct{}
	List  struct{}
	Get   struct{}
)

func (v *V1) Querier() *Query {
	return &Query{}
}

func (v *Query) List() *List {
	return &List{}
}

func (v *Query) Get() *Get {
	return &Get{}
}

type objectType string

const (
	typeVDC    objectType = "orgVdc"
	typeVAPP   objectType = "vApp"
	typeVM     objectType = "vm"
	typeEdgeGW objectType = "edgeGateway"
)

// getUUIDFromHref
func getUUIDFromHref(href string, idAtEnd bool) (string, error) {
	regex := `^https:\/\/.+([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`

	if idAtEnd {
		regex += `$`
	} else {
		regex += `.*$`
	}

	reGetID := regexp.MustCompile(regex)
	matchList := reGetID.FindAllStringSubmatch(href, -1)

	if len(matchList) == 0 {
		return "", fmt.Errorf("no match found")
	}
	return matchList[0][1], nil
}

// queryList
func queryList(objectType objectType) (govcd.Results, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		panic(err)
	}

	return c.Vmware.Query(map[string]string{
		"type": string(objectType),
	})
}

// queryListWithOptionalFilter
func queryListWithOptionalFilter(objectType objectType, filters map[string]string) (govcd.Results, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		panic(err)
	}

	filters["type"] = string(objectType)
	return c.Vmware.Query(filters)
}

// queryget
func queryGet(objectType objectType, name string) (govcd.Results, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		panic(err)
	}

	return c.Vmware.Query(map[string]string{
		"type":   string(objectType),
		"filter": "name==" + name,
	})
}

// queryGetWithOptionalFilter
func queryGetWithOptionalFilter(objectType objectType, name string, filters map[string]string) (govcd.Results, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		panic(err)
	}

	filters["type"] = string(objectType)
	return c.Vmware.Query(filters)
}

// VDC list all vdc informations
func (q *List) VDC() ([]*types.QueryResultOrgVdcRecordType, error) {
	r, err := queryList(typeVDC)
	return r.Results.OrgVdcRecord, err
}

// VDC get a vdc informations by name
func (q *Get) VDC(vdcName string) (*types.QueryResultOrgVdcRecordType, error) {
	r, err := queryGet(typeVDC, vdcName)
	return r.Results.OrgVdcRecord[0], err
}

// VAPP list all vapp informations
func (q *List) VAPP() ([]*types.QueryResultVAppRecordType, error) {
	r, err := queryList(typeVAPP)
	return r.Results.VAppRecord, err
}

// VAPP get a vapp informations by name
func (q *Get) VAPP(vappName string) (*types.QueryResultVAppRecordType, error) {
	r, err := queryGet(typeVAPP, vappName)
	return r.Results.VAppRecord[0], err
}

// VM list all vm informations
func (q *List) VM(vAppName string) ([]*types.QueryResultVMRecordType, error) {
	r, err := queryListWithOptionalFilter(typeVM, map[string]string{
		"filter": "containerName==" + vAppName,
	})

	for _, vm := range r.Results.VMRecord {
		id, err := getUUIDFromHref(vm.HREF, true)
		if err != nil {
			panic(err)
		}

		vm.ID = id
	}

	return r.Results.VMRecord, err
}

// VM get a vm informations by name
func (q *Get) VM(vmName, vAppName string) (*types.QueryResultVMRecordType, error) {
	r, err := queryGetWithOptionalFilter(typeVM, vmName, map[string]string{
		"filter": "containerName==" + vAppName,
		"name":   vmName,
	})
	if r.Results.VMRecord == nil {
		return nil, err
	}

	id, err := getUUIDFromHref(r.Results.VMRecord[0].HREF, true)
	if err == nil {
		r.Results.VMRecord[0].ID = id
	}

	return r.Results.VMRecord[0], err
}

// EdgeGW list all edgegw informations
func (q *List) EdgeGW() ([]*types.QueryResultEdgeGatewayRecordType, error) {
	r, err := queryList(typeEdgeGW)
	return r.Results.EdgeGatewayRecord, err
}

// EdgeGW get a edgegw informations by name
func (q *Get) EdgeGW(edgeGWName string) (*types.QueryResultEdgeGatewayRecordType, error) {
	r, err := queryGet(typeEdgeGW, edgeGWName)
	return r.Results.EdgeGatewayRecord[0], err
}
