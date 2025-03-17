/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package validators

import (
	"bytes"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// IPV4Range is a custom validator that checks if a string is a valid IPv4 range.
var IPV4Range = &CustomValidator{
	Key: "ipv4_range",
	Func: func(fl validator.FieldLevel) bool {
		// ipv4_range is a string in the form of "192.168.0.1-192.168.0.100"
		// We need to split the string into two parts and validate each part
		// as a valid IPv4 address.

		// Split the string into two parts
		parts := strings.Split(fl.Field().String(), "-")
		if len(parts) != 2 {
			return false
		}

		// Check if the first IP address is less than the second IP address
		firstIP := net.ParseIP(parts[0])
		secondIP := net.ParseIP(parts[1])
		if firstIP.To4() == nil || secondIP.To4() == nil {
			return false
		}

		if bytes.Compare(firstIP.To4(), secondIP.To4()) >= 0 {
			return false
		}

		return true
	},
}

// TCPUDPPort is a custom validator that checks if a string is a valid TCP or UDP port.
var TCPUDPPort = &CustomValidator{
	Key: "tcp_udp_port",
	Func: func(fl validator.FieldLevel) bool {
		var port int

		switch fl.Field().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			port = int(fl.Field().Int())
		case reflect.String:
			portStr := fl.Field().String()
			if portStr == "" {
				return false
			}

			// convert the string to an integer
			i, err := strconv.Atoi(fl.Field().String())
			if err != nil {
				return false
			}

			port = i

		}

		// check if the integer is a valid TCP or UDP port
		if port <= 0 || port > 65535 {
			return false
		}

		return true
	},
}

// TCPUDPPortRange is a custom validator that checks if a string is a valid TCP or UDP port range.
var TCPUDPPortRange = &CustomValidator{
	Key: "tcp_udp_port_range",
	Func: func(fl validator.FieldLevel) bool {
		// format of the string is "1-65535"
		// split the string into two parts
		p := strings.Split(fl.Field().String(), "-")

		if len(p) != 2 {
			return false
		}

		// convert the strings to integers
		start, err := strconv.Atoi(p[0])
		if err != nil {
			return false
		}

		end, err := strconv.Atoi(p[1])
		if err != nil {
			return false
		}

		// check if the integers are valid TCP or UDP ports
		if start <= 0 || start > 65535 {
			return false
		}

		if end <= 0 || end > 65535 {
			return false
		}

		// check if the start is less than the end
		if start >= end {
			return false
		}

		return true
	},
}
