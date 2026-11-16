/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package commoncloudavenue

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/internal/endpoints"
	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	caverrors "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

// ErrJobIDEmpty is returned when the API responds without a job ID, meaning
// the requested operation never started. Callers can match it with errors.Is
// to produce a tailored message.
var ErrJobIDEmpty = errors.New("job id is empty")

// MissingJobIDError builds the error to return when a write operation got an
// HTTP response with a non-error status code but no job ID in the body: the
// operation did not start and the raw response (status plus truncated body)
// is the only evidence of why the API declined it. It wraps ErrJobIDEmpty so
// callers can detect the case with errors.Is and keeps its wording in the
// message for substring matching.
func MissingJobIDError(r *resty.Response) error {
	return &apiCallError{
		statusCode: r.StatusCode(),
		message: fmt.Sprintf("%s: HTTP %s, response body: %s",
			ErrJobIDEmpty,
			r.Status(),
			caverrors.TruncateBody(strings.TrimSpace(r.String()), caverrors.MaxErrorBodyLen),
		),
		wrapped: ErrJobIDEmpty,
	}
}

// JobStatusMessage is a type for job status.
type JobStatusMessage string

const (
	// DONE is the done status message.
	DONE JobStatusMessage = "DONE"
	// FAILED is the failed status message.
	FAILED JobStatusMessage = "FAILED"
	// CREATED is the created status message.
	CREATED JobStatusMessage = "CREATED"
	// PENDING is the pending status message.
	PENDING JobStatusMessage = "PENDING"
	// INPROGRESS is the in progress status message.
	INPROGRESS JobStatusMessage = "IN_PROGRESS"
	// ERROR is the error status message.
	ERROR JobStatusMessage = "ERROR"
)

// JobCreatedAPIResponse - This is the response structure for the JobCreatedAPIResponse.
type JobCreatedAPIResponse struct {
	Message string `json:"message"`
	JobID   string `json:"jobId"`
}

// JobStatus - This is the response structure for the JobStatus.
type JobStatus struct {
	JobID   string `json:"jobId,omitempty"`
	Actions []struct {
		Name    string `json:"name"`
		Status  string `json:"status"`
		Details string `json:"details"`
	} `json:"actions"`
	Description string           `json:"description"`
	Name        string           `json:"name"`
	Status      JobStatusMessage `json:"status"`
}

// GetJobStatus - Returns the status of a job.
func (j *JobCreatedAPIResponse) GetJobStatus() (response *JobStatus, err error) {
	response = new(JobStatus)
	response.JobID = j.JobID
	if err := response.Refresh(); err != nil {
		return nil, err
	}

	return response, response.Refresh()
}

// Refresh - Refreshes the job status.
func (j *JobStatus) Refresh() error {
	if j.JobID == "" {
		return fmt.Errorf("cannot refresh job status: %w", ErrJobIDEmpty)
	}

	jobID := j.JobID

	c := clientcloudavenue.GetClient()

	r, err := c.R().
		SetResult(&[]JobStatus{}).
		SetError(&APIErrorResponse{}).
		SetPathParams(map[string]string{
			"jobId": j.JobID,
		}).
		Get(endpoints.JobStatusGet)
	if err != nil {
		return err
	}

	if r.IsError() {
		return fmt.Errorf("error on API call: %w", ToError(r))
	}

	x := *r.Result().(*[]JobStatus)
	if len(x) == 0 {
		return fmt.Errorf("job with ID %s not found", jobID)
	}
	*j = x[0]

	j.JobID = jobID

	return nil
}

// IsDone - Returns true if the job is done with a success.
func (j *JobStatus) IsDone() bool {
	return j.Status == DONE
}

// OnError - Returns true if the job is done with an error.
func (j *JobStatus) OnError() bool {
	return j.Status == FAILED || j.Status == ERROR
}

// Wait - Waits for the job to be done
// refreshInterval - The interval in seconds between each refresh
// timeout - The timeout in seconds.
func (j *JobStatus) Wait(refreshInterval, timeout int) error {
	err := j.Refresh()
	if err != nil {
		return err
	}

	if j.IsDone() {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	return j.WaitWithContext(ctx, refreshInterval)
}

// WaitWithContext - Waits for the job to be done
// refreshInterval - The interval in seconds between each refresh.
func (j *JobStatus) WaitWithContext(ctx context.Context, refreshInterval int) error {
	if _, deadlineSet := ctx.Deadline(); !deadlineSet {
		return j.Wait(refreshInterval, 5*60)
	}

	err := j.Refresh()
	if err != nil {
		return err
	}

	if j.IsDone() {
		return nil
	}

	ticker := time.NewTicker(time.Duration(refreshInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout reached (%w)", ctx.Err())
		case <-ticker.C:
			err := j.Refresh()
			if err != nil {
				return err
			}

			if j.IsDone() {
				return nil
			}

			if j.OnError() {
				// find the first action that failed
				for _, a := range j.Actions {
					if a.Status == string(FAILED) || a.Status == string(ERROR) {
						return fmt.Errorf("job failed: %s", a.Details)
					}
				}

				return fmt.Errorf("job failed: %s", j.Description)
			}
		}
	}
}
