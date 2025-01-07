/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package urn

import "regexp"

// ExtractUUID finds an UUID in the input string
// Returns an empty string if no UUID was found.
func ExtractUUID(input string) string {
	reGetID := regexp.MustCompile(`([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)
	matchListIDs := reGetID.FindAllStringSubmatch(input, -1)
	if len(matchListIDs) > 0 && len(matchListIDs[0]) > 0 {
		return matchListIDs[len(matchListIDs)-1][len(matchListIDs[0])-1]
	}
	return ""
}

// IsValid returns true if the URN is valid.
// Checks if the URN is of the specified type and if it is a valid UUIDv4.
func IsValid(urn string) bool {
	if len(urn) == 0 {
		return false
	}

	u := URN(urn)

	for _, prefix := range URNs {
		// Check if the URN is of the specified type.
		if u.IsType(prefix) {
			// Check if the URN is a valid UUIDv4.
			return isUUIDV4(extractUUIDv4(urn, prefix))
		}
	}
	return false
}

// Normalize returns the URN with the prefix if prefix is missing.
func Normalize(prefix URN, uuid string) URN {
	u := URN(uuid)
	if u.ContainsPrefix() {
		return u
	}

	if prefix.isEmpty() {
		return ""
	}

	return prefix + u
}
