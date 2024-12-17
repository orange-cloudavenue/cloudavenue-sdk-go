## 0.18.0 (Unreleased)

### :rocket: **New Features**

* `edgeGateway` - Add funcs to manage a security group `CreateFirewallSecurityGroup` and `GetFirewallSecurityGroup`. (GH-159)
* `vdcGroup` - Add funcs to manage a security group `CreateFirewallSecurityGroup` and `GetFirewallSecurityGroup`. (GH-159)

## 0.17.0 (December 16, 2024)
### :rotating_light: **Breaking Changes**

* `client` - Now the `endpoint` parameter for client opts is removed. Use `url` instead. (GH-168)

### :dependabot: **Dependencies**

* deps: bumps github.com/vmware/go-vcloud-director/v2 from 2.26.0 to 2.26.1 (GH-158)
* deps: bumps golang.org/x/sync from 0.9.0 to 0.10.0 (GH-157)

## 0.16.0 (December 13, 2024)

### :rocket: **New Features**

* `vdcGroup` - Add new functions to manage Network Isolated in VDC Group. (In VDCGroup object `CreateNetworkIsolated`, `GetNetworkIsolated`, `Update`, `Delete`) (GH-161)
* `vdc` - Add new functions to manage Network Isolated in VDC. (In VDC object `CreateNetworkIsolated`, `GetNetworkIsolated`, `Update`, `Delete`) (GH-161)

## 0.15.1 (December 12, 2024)
## 0.15.0 (December 12, 2024)

### :rocket: **New Features**

* `edgegateway` - Add new functions to manage ALB Virtual Service (Advanced Load Balancer). (In EdgeGateway object `CreateALBVirtualService`, `GetALBVirtualService`, `Update`, `Delete`) (GH-160)
* `vdcGroup` - Add new fuctions to manage Firewall in VDC Group. (In VDCGroup object `CreateFirewall`, `GetFirewall`. In VdcGroupFirewall Object `GetRules`, `IsEnabled`, `UpdateFirewall`, `Delete`) (GH-153)

### :dependabot: **Dependencies**

* deps: bumps github.com/go-resty/resty/v2 from 2.16.0 to 2.16.2 (GH-150)

## 0.14.0 (December  2, 2024)
### :rotating_light: **Breaking Changes**

* `uuid` - The `uuid` package has been renamed to `urn` to better reflect its purpose. Replace your imports from `"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/uuid"` to `"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"` to avoid breaking changes. (GH-154)

## 0.13.0 (November 26, 2024)
### :rotating_light: **Breaking Changes**

* The struct `EdgeGw` has renamed to `EdgeGatewayType`. (GH-152)
### :information_source: **Notes**

* `edgateway/GetByID` - Now the function edgateway/GetByID is deprecated and will be removed in the next version. Use the function edgateway/Get instead. (GH-152)
* `edgateway/GetByName` - Now the function edgateway/GetByName is deprecated and will be removed in the next version. Use the function edgateway/Get instead. (GH-152)

## 0.12.7 (November 25, 2024)

### :rocket: **New Features**

* `client` - Add `adminVDC` in the client. (GH-149)

## 0.12.6 (November 22, 2024)

### :rocket: **New Features**

* `client` - Add `org` and `adminorg` in the client. (GH-148)
* `v1/vdc` - Add `GetVDC`/ `GetVDCGroup`/ `GetVDCOrVDCGroup` (GH-148)

## 0.12.5 (November 22, 2024)
## 0.12.4 (November 22, 2024)

### :bug: **Bug Fixes**

* `netbackup` - Fixed a bug where Netbackup API URL changed. (GH-146)

### :dependabot: **Dependencies**

* deps: bumps github.com/go-resty/resty/v2 from 2.15.3 to 2.16.0 (GH-143)
* deps: bumps github.com/hashicorp/terraform-plugin-sdk/v2 from 2.34.0 to 2.35.0 (GH-140)
* deps: bumps golang.org/x/mod from 0.21.0 to 0.22.0 (GH-142)
* deps: bumps golang.org/x/sync from 0.8.0 to 0.9.0 (GH-141)

## 0.12.3 (November 19, 2024)

### :bug: **Bug Fixes**

* `edge_gateway` - Fix a bug where the remaining bandwidth was not being calculated correctly. (GH-144)

## 0.12.2 (November  3, 2024)

### :dependabot: **Dependencies**

* deps: bumps github.com/go-resty/resty/v2 from 2.15.0 to 2.15.3 (GH-139)
* deps: bumps github.com/vmware/go-vcloud-director/v2 from 2.25.0 to 2.26.0 (GH-137)

## 0.12.1 (September 18, 2024)

### :tada: **Improvements**

* `BMS` - Add func GetBMSByHostname (GH-135)

## 0.12.0 (September 17, 2024)

### :rocket: **New Features**

* `cav/bms` - Add support for `bms` endpoints (GH-130)

### :dependabot: **Dependencies**

* deps: bumps github.com/go-resty/resty/v2 from 2.14.0 to 2.15.0 (GH-132)
* deps: bumps golang.org/x/mod from 0.20.0 to 0.21.0 (GH-129)

## 0.11.1 (August 23, 2024)

### :bug: **Bug Fixes**

* `vdc_rules` - Fix bad regex on storage profile custom name validation. (GH-128)

### :dependabot: **Dependencies**

* deps: bumps github.com/go-resty/resty/v2 from 2.13.1 to 2.14.0 (GH-125)
* deps: bumps github.com/vmware/go-vcloud-director/v2 from 2.23.0 to 2.25.0 (GH-123)
* deps: bumps golang.org/x/mod from 0.17.0 to 0.20.0 (GH-126)
* deps: bumps golang.org/x/sync from 0.7.0 to 0.8.0 (GH-127)

## 0.11.0 (July 31, 2024)

### :rotating_light: **Breaking Changes**

* `v1/vdc` - Now the `New`/`Update`/`Delete` methods of the `vdc` package require a `context.Context` as the first argument. (GH-117)

### :rocket: **New Features**

* `job` - Add `WaitWithContext` method to `job` package. This method allows users to wait for a job to complete with a context (Context with timeout). If context does not have a timeout, it will wait 90 seconds by default. (GH-117)

### :bug: **Bug Fixes**

* `job` - Fix nil pointer dereference in `wait` command when the job is not found. (GH-67)

### :information_source: **Notes**

* `job` - Now if the job return a error (ERROR or FAILED status code), the `Wait` method will return the error. (GH-116)

### :dependabot: **Dependencies**

* deps: bumps actions/download-artifact from 4.1.4 to 4.1.8 (GH-122)
* deps: bumps dependabot/fetch-metadata from 2.0.0 to 2.2.0 (GH-120)
* deps: bumps github.com/aws/aws-sdk-go from 1.51.21 to 1.55.5 (GH-124)
* deps: bumps github.com/go-resty/resty/v2 from 2.12.0 to 2.13.1 (GH-113)
* deps: bumps github.com/hashicorp/terraform-plugin-sdk/v2 from 2.33.0 to 2.34.0 (GH-114)
* deps: bumps github.com/sethvargo/go-envconfig from 1.0.1 to 1.1.0 (GH-119)
* deps: bumps golangci/golangci-lint-action from 4 to 6 (GH-121)

## 0.10.1 (April 26, 2024)

### :bug: **Bug Fixes**

* `querier/get` - Fix nil pointer when resource doesn't exist. (GH-101)

### :dependabot: **Dependencies**

* deps: bumps dependabot/fetch-metadata from 1.6.0 to 2.0.0 (GH-98)
* deps: bumps github.com/aws/aws-sdk-go from 1.51.0 to 1.51.1 (GH-96)
* deps: bumps github.com/aws/aws-sdk-go from 1.51.1 to 1.51.21 (GH-107)
* deps: bumps github.com/go-resty/resty/v2 from 2.11.0 to 2.12.0 (GH-97)
* deps: bumps github.com/sethvargo/go-envconfig from 0.9.0 to 1.0.1 (GH-103)
* deps: bumps github.com/vmware/go-vcloud-director/v2 from 2.22.0 to 2.23.0 (GH-99)
* deps: bumps golang.org/x/mod from 0.16.0 to 0.17.0 (GH-102)
* deps: bumps golang.org/x/sync from 0.6.0 to 0.7.0 (GH-106)

## 0.10.0 (March 15, 2024)

### :rocket: **New Features**

* `edgegateway` - Support Dedicated EdgeGateway Medium and Large. (GH-90)

### :dependabot: **Dependencies**

* deps: bumps actions/download-artifact from 4.1.1 to 4.1.2 (GH-85)
* deps: bumps actions/download-artifact from 4.1.2 to 4.1.4 (GH-91)
* deps: bumps github.com/aws/aws-sdk-go from 1.50.10 to 1.51.0 (GH-95)
* deps: bumps github.com/hashicorp/terraform-plugin-sdk/v2 from 2.32.0 to 2.33.0 (GH-89)
* deps: bumps golang.org/x/mod from 0.14.0 to 0.16.0 (GH-94)
* deps: bumps golangci/golangci-lint-action from 3.7.0 to 4.0.0 (GH-86)
* deps: bumps softprops/action-gh-release from 1 to 2 (GH-92)

## 0.9.1 (February  6, 2024)

### :bug: **Bug Fixes**

* `rules/vdc` - Fix custom storage profile class. (GH-83)

### :dependabot: **Dependencies**

* deps: bumps github.com/aws/aws-sdk-go from 1.50.7 to 1.50.10 (GH-81)
* deps: bumps github.com/hashicorp/terraform-plugin-sdk/v2 from 2.31.0 to 2.32.0 (GH-82)

## 0.9.0 (January 31, 2024)

### :rocket: **New Features**

* `rules/vdc` - Now support custom storage profile class. (GH-78)

### :dependabot: **Dependencies**

* deps: bumps github.com/aws/aws-sdk-go from 1.49.16 to 1.50.7 (GH-80)
* deps: bumps github.com/vmware/go-vcloud-director/v2 from 2.21.0 to 2.22.0 (GH-54)

## 0.8.1 (January 30, 2024)

### :bug: **Bug Fixes**

* `client` - Fix bug `nil` pointer dereference in `client` package if no `opts` are passed to `New()` func. (GH-75)

## 0.8.0 (January 29, 2024)

### :rotating_light: **Breaking Changes**

* `console` - Remove funcs `S3IsEnabled()`, `GetS3Endpoint()`, `IsVCDAEnabled()` and `GetVCDAEndpoint()`. (GH-73)

### :rocket: **New Features**

* `client/cloudavenue` - Add func `GetURL()` to get the cloudavenue url. (GH-73)
* `client/s3` - Now the `S3` client return an error if the s3 service is not available in the location. (GH-73)
* `client` - Add Validation for the creation of a new client (CloudAvenue and Netbackup). (GH-73)
* `console/service` - Add funcs `IsEnabled()` and `GetEndpoint()`. (GH-73)
* `console` - Add func `CheckOrganizationName()` to check if the organization name is valid without creating new client. (GH-73)
* `console`- Add func `Services()` to get the services available in the console. (GH-73)
* `consoles` - Add consoles `console5`, `console7`, `console8` and `console9`. (GH-73)
* `errors` - Add `errors` package. The following errors are available: `ErrNotFound`, `ErrEmpty` and `ErrInvalidFormat` (GH-73)
* `v1/Vmware` - Now the `V1()` function returns a Vmware object. (GH-73)

### :information_source: **Notes**

* `client` - The `CLOUDAVENUE_ENDPOINT` environment variable has been renamed to `CLOUDAVENUE_URL`. (GH-73)
* `client` - The `NETBACKUP_ENDPOINT` environment variable has been renamed to `NETBACKUP_URL`. (GH-73)

### :dependabot: **Dependencies**

* deps: bumps actions/download-artifact from 4.1.0 to 4.1.1 (GH-70)

## 0.7.1 (January 18, 2024)

### :rocket: **New Features**

* `publicIP/GetIPByJob` - Retrieve the public IP address by a Job. (GH-71)

## 0.7.0 (January  9, 2024)

### :rocket: **New Features**

* `client` - Improve seep of `client`. Use cached if already connected. (GH-63)
* `querier` - Add querier vmware. Permit to Get/List VDC/VAPP/VM/EdgeGateway. (GH-65)
* `vdc/rules` - Add storage profile rules for VDC. (GH-64)

## 0.6.1 (January  8, 2024)

### :dependabot: **Dependencies**

* deps: bumps actions/download-artifact from 3.0.2 to 4.1.0 (GH-61)
* deps: bumps actions/setup-go from 4 to 5 (GH-50)
* deps: bumps actions/upload-artifact from 3 to 4 (GH-53)
* deps: bumps github.com/aws/aws-sdk-go from 1.47.9 to 1.49.16 (GH-62)
* deps: bumps github.com/go-resty/resty/v2 from 2.10.0 to 2.11.0 (GH-56)
* deps: bumps github.com/hashicorp/terraform-plugin-sdk/v2 from 2.30.0 to 2.31.0 (GH-55)
* deps: bumps github/codeql-action from 2 to 3 (GH-52)
* deps: bumps golang.org/x/sync from 0.1.0 to 0.6.0 (GH-60)

## 0.6.0 (December  4, 2023)

### :rotating_light: **Breaking Changes**

* `edgegateway` - Change attribute name to respect naming ToService/VlanID to ToService/VLANID (GH-45)
* `edgegateway` - Change funcs name to respect naming rules from `GetVlanID` to `GetVLANID` (GH-45)
* `netbackup` - Change funcs name to respect naming rules from `GetVdcByID` to `GetVDCByID` (GH-45)
* `netbackup` - Change funcs name to respect naming rules from `GetVdcByIdentifier` to `GetVDCByIdentifier` (GH-45)
* `netbackup` - Change funcs name to respect naming rules from `GetVdcByNameOrIdentifier` to `GetVDCByNameOrIdentifier` (GH-45)
* `netbackup` - Change funcs name to respect naming rules from `GetVdcByName` to `GetVDCByName` (GH-45)
* `vdc` - Change funcs name from `GetVcpuInMhz2` to `GetVCPUInMhz` (GH-45)

### :rocket: **New Features**

* `vdc` - Add `SetStorageProfiles`, `SetVCPUInMhz2` and `Set` funcs (GH-43)
* `vdc` - Add vdc rules management (GH-45)

### :tada: **Improvements**

* `Lint` - Add lint for upper case var-naming rules (GH-45)

### :bug: **Bug Fixes**

* `publicip` - Fix GetIP return now the good public IP. (GH-40)
* `publicip` - Fix GetJobStatus return now the good status. (GH-40)

### :information_source: **Notes**

* `vdc` - Refactor `vdc` to use `infrapi` and `vmware` packages (GH-42)

## 0.5.5 (November 20, 2023)

### :information_source: **Notes**

* `netbackup/` - Reorganize the NetBackup files into a directory. (GH-39)

## 0.5.4 (November 20, 2023)

### :bug: **Bug Fixes**

* `client/cloudavenue` - permit to configure the client only by environment variables. (GH-38)
* `client/s3` - Now the client use cloudavenue settings after environment variables is evaluated. (GH-38)

## 0.5.3 (November 17, 2023)

## 0.5.2 (November 16, 2023)

### :bug: **Bug Fixes**

* `s3/login` - Fix bug where `s3/login` would fail if the user does not have a system key. (GH-31)

## 0.5.1 (November 15, 2023)

### :bug: **Bug Fixes**

* `s3/GetCanonicalID` - Fix GetCanonicalID to return the correct canonical ID for the account. (GH-30)

## 0.5.0 (November 15, 2023)

## 0.5.0 (November 15, 2023)

### :tada: **Improvements**

* `vcda/ip` - Allow to List/Get/Create/Delete VCDA IP (GH-28)

## 0.4.1 (November 13, 2023)

### :rocket: **New Features**

* `s3/Sync - Force synchronization of a bucket ([GH-19](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/19))

### :tada: **Improvements**

* `s3/NewCredential` - Remove `username` parameter. ([GH-21](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/21))

### :dependabot: **Dependencies**

* deps: bumps github.com/aws/aws-sdk-go from 1.45.26 to 1.47.9 ([GH-23](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/23))

## 0.4.0 (November  7, 2023)

### :tada: **Improvements**

* `s3/credential` - Allow to List/Get/Delete OSE User Credential ([GH-18](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/18))
* `s3/user` - Allow to List/Get OSE User and Get Canonical ID ([GH-18](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/18))

## 0.3.1 (October 24, 2023)

## 0.3.0 (October 17, 2023)

## 0.2.0 (October 16, 2023)

## 0.1.0 (October 16, 2023)

### :tada: **Improvements**

* `v1/edgegw` - Add `GetAllowedBandwidthValues` function to get the Allowed Bandwidth Values of the Edge Gateway in Mbps. ([GH-9](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/9))
* `v1/edgegw` - Add `GetBandwidthCapacityRemaining` function to get the Bandwidth Capacity Remaining of the Edge Gateway in Mbps. ([GH-9](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/9))
* `v1/t0` - Add `GetBandwidthCapacity` function to get the Bandwidth Capacity of the T0 in Mbps. ([GH-9](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/9))

### :dependabot: **Dependencies**

* deps: bumps github.com/go-resty/resty/v2 from 2.9.1 to 2.10.0 ([GH-13](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/13))

## 0.0.3 (October  9, 2023)

## 0.0.2 (October  9, 2023)
