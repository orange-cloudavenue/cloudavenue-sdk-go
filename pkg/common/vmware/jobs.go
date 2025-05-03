/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package commonvmware

import (
	"github.com/go-resty/resty/v2"

	"github.com/vmware/go-vcloud-director/v2/govcd"
)

func NewAndWait(c *govcd.Client, r *resty.Response) error {
	task := govcd.NewTask(c)
	// the task is not in the response, so we need to get it from the header
	task.Task.HREF = r.Header().Get("Location")

	return task.WaitTaskCompletion()
}
