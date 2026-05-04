/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgeloadbalancer

// Shared test constants used across all test files in this package.
const (
	// Virtual service fixtures.
	testSuccess             = "success"
	testVirtualServiceName1 = "virtualServiceName1"
	testVirtualServiceDesc1 = "virtualServiceDescription1"
	testVirtualServiceName2 = "virtualServiceName2"
	testVirtualServiceDesc2 = "virtualServiceDescription2"
	testRuleName            = "ruleName"
	testRuleName2           = "ruleName2"
	testRuleName3           = "ruleName3"

	// Network fixtures.
	testIPAddress = "192.168.0.1"
	testIPSingle  = "12.23.34.45"
	testIPCIDR    = "12.23.34.0/24"
	testIPRange   = "12.23.34.0-12.23.34.100"

	// HTTP path fixtures.
	testPath1   = "/path1"
	testPath2   = "/path2"
	testNewPath = "/newpath"

	// HTTP query param fixtures.
	testQuery1 = "key1=value1"
	testQuery2 = "key2=value2"

	// HTTP header fixtures.
	testHeaderUserAgent       = "User-Agent"
	testHeaderAccept          = "Accept"
	testHeaderXForwardedFor   = "X-Forwarded-For"
	testHeaderXForwardedProto = "X-Forwarded-Proto"
	testHeaderXCustom         = "X-Custom-Header"

	// HTTP header value fixtures.
	testHeaderValueMozilla = "Mozilla/5.0"
	testHeaderValueCurl    = "curl/7.64.1"
	testHeaderValueTest    = "test"
	testHeaderValue1       = "value1"
	testHeaderValue2       = "value2"

	// Cookie fixtures.
	testCookieName  = "session_id"
	testCookieValue = "abc123"

	// Domain/redirect fixtures.
	testDomain       = "example.com"
	testRedirectCode = "301-303"

	// Error state fixtures.
	testErrorRefresh         = "error-refresh"
	testErrorVSValidation    = "error-virtualserviceValidation"
	testErrorGetVS           = "error-getVirtualService"
	testErrorValidationModel = "error-validation-model"
	testErrorDelete          = "error-delete"
	testErrorRefreshShort    = "refresh-error"
	testErrorGetShort        = "error-get"
	testErrorValidation      = "error-validation"
	testErrorGetAllCerts     = "error-get-all-certificates"

	// Pool fixtures.
	testPoolName1            = "pool1"
	testPoolName1Desc        = "pool1 description"
	testPoolName2            = "pool2"
	testPoolName2Desc        = "pool2 description"
	testPoolMonitorHTTP      = "monitor HTTP"
	testPoolMonitorTCP       = "monitor TCP"
	testPoolPersistence      = "persistence profile"
	testPoolMembersStatus    = "All members are up"
	testPoolPoule2           = "poule 2"
	testPoolPouleDesc        = "poule description"
	testPoolParamEdgeEmpty   = "param-edgeGatewayID-empty"
	testPoolParamEdgeInvalid = "param-edgeGatewayID-invalid-id"

	// Service engine group fixtures.
	testEdgeName = "edge_name"

	// Shared name fixture.
	testName = "name"

	// Security policy fixtures.
	testJSONBody   = "{\"key\":\"value\"}"
	testBase64Body = "eyJrZXkiOiJ2YWx1ZSJ9Cg=="

	// Content type fixtures.
	testContentTypeJSON = "application/json"
	testContentTypeHTML = "text/html"
)
