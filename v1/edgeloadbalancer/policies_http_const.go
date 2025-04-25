/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgeloadbalancer

const (
	PoliciesHTTPProtocolHTTP  PoliciesHTTPProtocol = "HTTP"
	PoliciesHTTPProtocolHTTPS PoliciesHTTPProtocol = "HTTPS"

	PoliciesHTTPMatchCriteriaCriteriaISIN              PoliciesHTTPMatchCriteriaCriteria = "IS_IN"
	PoliciesHTTPMatchCriteriaCriteriaISNOTIN           PoliciesHTTPMatchCriteriaCriteria = "IS_NOT_IN"
	PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH        PoliciesHTTPMatchCriteriaCriteria = "BEGINS_WITH"
	PoliciesHTTPMatchCriteriaCriteriaDOESNOTBEGINWITH  PoliciesHTTPMatchCriteriaCriteria = "DOES_NOT_BEGIN_WITH"
	PoliciesHTTPMatchCriteriaCriteriaCONTAINS          PoliciesHTTPMatchCriteriaCriteria = "CONTAINS"
	PoliciesHTTPMatchCriteriaCriteriaDOESNOTCONTAIN    PoliciesHTTPMatchCriteriaCriteria = "DOES_NOT_CONTAIN"
	PoliciesHTTPMatchCriteriaCriteriaENDSWITH          PoliciesHTTPMatchCriteriaCriteria = "ENDS_WITH"
	PoliciesHTTPMatchCriteriaCriteriaDOESNOTENDWITH    PoliciesHTTPMatchCriteriaCriteria = "DOES_NOT_END_WITH"
	PoliciesHTTPMatchCriteriaCriteriaEQUALS            PoliciesHTTPMatchCriteriaCriteria = "EQUALS"
	PoliciesHTTPMatchCriteriaCriteriaDOESNOTEQUAL      PoliciesHTTPMatchCriteriaCriteria = "DOES_NOT_EQUAL"
	PoliciesHTTPMatchCriteriaCriteriaREGEXMATCH        PoliciesHTTPMatchCriteriaCriteria = "REGEX_MATCH"
	PoliciesHTTPMatchCriteriaCriteriaREGEXDOESNOTMATCH PoliciesHTTPMatchCriteriaCriteria = "REGEX_DOES_NOT_MATCH"
	PoliciesHTTPMatchCriteriaCriteriaEXISTS            PoliciesHTTPMatchCriteriaCriteria = "EXISTS"
	PoliciesHTTPMatchCriteriaCriteriaDOESNOTEXIST      PoliciesHTTPMatchCriteriaCriteria = "DOES_NOT_EXIST"

	PoliciesHTTPActionHeaderRewriteActionADD     PoliciesHTTPActionHeaderRewriteAction = "ADD"
	PoliciesHTTPActionHeaderRewriteActionREMOVE  PoliciesHTTPActionHeaderRewriteAction = "REMOVE"
	PoliciesHTTPActionHeaderRewriteActionREPLACE PoliciesHTTPActionHeaderRewriteAction = "REPLACE"

	PoliciesHTTPMethodGET       PoliciesHTTPMethod = "GET"
	PoliciesHTTPMethodPOST      PoliciesHTTPMethod = "POST"
	PoliciesHTTPMethodPUT       PoliciesHTTPMethod = "PUT"
	PoliciesHTTPMethodDELETE    PoliciesHTTPMethod = "DELETE"
	PoliciesHTTPMethodPATCH     PoliciesHTTPMethod = "PATCH"
	PoliciesHTTPMethodOPTIONS   PoliciesHTTPMethod = "OPTIONS"
	PoliciesHTTPMethodTRACE     PoliciesHTTPMethod = "TRACE"
	PoliciesHTTPMethodCONNECT   PoliciesHTTPMethod = "CONNECT"
	PoliciesHTTPMethodPROPFIND  PoliciesHTTPMethod = "PROPFIND"
	PoliciesHTTPMethodPROPPATCH PoliciesHTTPMethod = "PROPPATCH"
	PoliciesHTTPMethodMKCOL     PoliciesHTTPMethod = "MKCOL"
	PoliciesHTTPMethodCOPY      PoliciesHTTPMethod = "COPY"
	PoliciesHTTPMethodMOVE      PoliciesHTTPMethod = "MOVE"
	PoliciesHTTPMethodLOCK      PoliciesHTTPMethod = "LOCK"
	PoliciesHTTPMethodUNLOCK    PoliciesHTTPMethod = "UNLOCK"

	PoliciesHTTPConnectionActionALLOW PoliciesHTTPConnectionAction = "ALLOW"
	PoliciesHTTPConnectionActionCLOSE PoliciesHTTPConnectionAction = "CLOSE"
)

// * Var HTTP.
var (
	PoliciesHTTPProtocols = []PoliciesHTTPProtocol{
		PoliciesHTTPProtocolHTTP,
		PoliciesHTTPProtocolHTTPS,
	}

	PoliciesHTTPProtocolsString = sliceAnyToSliceString(PoliciesHTTPProtocols)

	PoliciesHTTPMethodsMatch = []PoliciesHTTPMethod{
		PoliciesHTTPMethodGET,
		PoliciesHTTPMethodPOST,
		PoliciesHTTPMethodPUT,
		PoliciesHTTPMethodDELETE,
		PoliciesHTTPMethodPATCH,
		PoliciesHTTPMethodOPTIONS,
		PoliciesHTTPMethodTRACE,
		PoliciesHTTPMethodCONNECT,
		PoliciesHTTPMethodPROPFIND,
		PoliciesHTTPMethodPROPPATCH,
		PoliciesHTTPMethodMKCOL,
		PoliciesHTTPMethodCOPY,
		PoliciesHTTPMethodMOVE,
		PoliciesHTTPMethodLOCK,
		PoliciesHTTPMethodUNLOCK,
	}
	PoliciesHTTPMethodsMatchString = sliceAnyToSliceString(PoliciesHTTPMethodsMatch)
)

// * Var Action.
var (
	PoliciesHTTPActionHeaderRewriteActions = []PoliciesHTTPActionHeaderRewriteAction{
		PoliciesHTTPActionHeaderRewriteActionADD,
		PoliciesHTTPActionHeaderRewriteActionREMOVE,
		PoliciesHTTPActionHeaderRewriteActionREPLACE,
	}
	PoliciesHTTPActionHeaderRewriteActionsString = sliceAnyToSliceString(PoliciesHTTPActionHeaderRewriteActions)
)

// * Var for match criteria.
var (
	PoliciesHTTPMatchCriteriaProtocols       = PoliciesHTTPProtocols
	PoliciesHTTPMatchCriteriaProtocolsString = PoliciesHTTPProtocolsString

	PoliciesHTTPMethodMatchCriteria = []PoliciesHTTPMatchCriteriaCriteria{
		PoliciesHTTPMatchCriteriaCriteriaISIN,
		PoliciesHTTPMatchCriteriaCriteriaISNOTIN,
	}
	PoliciesHTTPMethodMatchCriteriaString = sliceAnyToSliceString(PoliciesHTTPMethodMatchCriteria)

	PoliciesHTTPServicePortMatchCriteria       = PoliciesHTTPMethodMatchCriteria
	PoliciesHTTPServicePortMatchCriteriaString = PoliciesHTTPMethodMatchCriteriaString

	PoliciesHTTPClientIPMatchCriteria       = PoliciesHTTPMethodMatchCriteria
	PoliciesHTTPClientIPMatchCriteriaString = PoliciesHTTPMethodMatchCriteriaString

	PoliciesHTTPStatusCodeMatchCriteria       = PoliciesHTTPMethodMatchCriteria
	PoliciesHTTPStatusCodeMatchCriteriaString = PoliciesHTTPMethodMatchCriteriaString

	PoliciesHTTPPathMatchCriteria = []PoliciesHTTPMatchCriteriaCriteria{
		PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTBEGINWITH,
		PoliciesHTTPMatchCriteriaCriteriaCONTAINS,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTCONTAIN,
		PoliciesHTTPMatchCriteriaCriteriaENDSWITH,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTENDWITH,
		PoliciesHTTPMatchCriteriaCriteriaEQUALS,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTEQUAL,
		PoliciesHTTPMatchCriteriaCriteriaREGEXMATCH,
		PoliciesHTTPMatchCriteriaCriteriaREGEXDOESNOTMATCH,
	}

	PoliciesHTTPPathMatchCriteriaString = sliceAnyToSliceString(PoliciesHTTPPathMatchCriteria)

	PoliciesHTTPLocationMatchCriteria       = PoliciesHTTPPathMatchCriteria
	PoliciesHTTPLocationMatchCriteriaString = PoliciesHTTPPathMatchCriteriaString

	PoliciesHTTPHeaderMatchCriteria = []PoliciesHTTPMatchCriteriaCriteria{
		PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTBEGINWITH,
		PoliciesHTTPMatchCriteriaCriteriaCONTAINS,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTCONTAIN,
		PoliciesHTTPMatchCriteriaCriteriaENDSWITH,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTENDWITH,
		PoliciesHTTPMatchCriteriaCriteriaEQUALS,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTEQUAL,
		PoliciesHTTPMatchCriteriaCriteriaEXISTS,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTEXIST,
	}

	PoliciesHTTPHeaderMatchCriteriaString = sliceAnyToSliceString(PoliciesHTTPHeaderMatchCriteria)

	PoliciesHTTPCookieMatchCriteria = []PoliciesHTTPMatchCriteriaCriteria{
		PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTBEGINWITH,
		PoliciesHTTPMatchCriteriaCriteriaCONTAINS,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTCONTAIN,
		PoliciesHTTPMatchCriteriaCriteriaENDSWITH,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTENDWITH,
		PoliciesHTTPMatchCriteriaCriteriaEQUALS,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTEQUAL,
		PoliciesHTTPMatchCriteriaCriteriaEXISTS,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTEXIST,
	}
	PoliciesHTTPCookieMatchCriteriaString = sliceAnyToSliceString(PoliciesHTTPCookieMatchCriteria)
)

// * Var for Action.
var (
	PoliciesHTTPRedirectActionProtocols       = PoliciesHTTPProtocols
	PoliciesHTTPRedirectActionProtocolsString = PoliciesHTTPProtocolsString

	PoliciesHTTPRedirectActionStatusCodes       = []int{301, 302, 307}
	PoliciesHTTPRedirectActionStatusCodesString = sliceAnyToSliceString(PoliciesHTTPRedirectActionStatusCodes)
)

// * Var Response.
var (
	PoliciesHTTPResponseLocationHeaderMatchCriteria = []PoliciesHTTPMatchCriteriaCriteria{
		PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTBEGINWITH,
		PoliciesHTTPMatchCriteriaCriteriaCONTAINS,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTCONTAIN,
		PoliciesHTTPMatchCriteriaCriteriaENDSWITH,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTENDWITH,
		PoliciesHTTPMatchCriteriaCriteriaEQUALS,
		PoliciesHTTPMatchCriteriaCriteriaDOESNOTEQUAL,
		PoliciesHTTPMatchCriteriaCriteriaREGEXMATCH,
		PoliciesHTTPMatchCriteriaCriteriaREGEXDOESNOTMATCH,
	}
	PoliciesHTTPResponseLocationHeaderMatchCriteriaString = sliceAnyToSliceString(PoliciesHTTPResponseLocationHeaderMatchCriteria)
)

// * Var for Connection.
var (
	PoliciesHTTPConnectionActions = []PoliciesHTTPConnectionAction{
		PoliciesHTTPConnectionActionALLOW,
		PoliciesHTTPConnectionActionCLOSE,
	}
	PoliciesHTTPConnectionActionsString = sliceAnyToSliceString(PoliciesHTTPConnectionActions)
)

// * Var for ActionSendReponseStatusCode.
var (
	PoliciesHTTPActionResponseStatusCodes       = []int64{200, 204, 403, 404, 429, 501}
	PoliciesHTTPActionResponseStatusCodesString = sliceAnyToSliceString(PoliciesHTTPActionResponseStatusCodes)
)

// * Var for content type.
var (
	PoliciesHTTPActionContentTypes = []string{
		"application/json",
		"text/plain",
		"text/html",
	}
	PoliciesHTTPActionContentTypesString = sliceAnyToSliceString(PoliciesHTTPActionContentTypes)
)

// * Var for RedirectStatusCode.
var (
	PoliciesHTTPRedirectStatusCodes       = []int64{301, 302, 307}
	PoliciesHTTPRedirectStatusCodesString = sliceAnyToSliceString(PoliciesHTTPRedirectStatusCodes)
)
