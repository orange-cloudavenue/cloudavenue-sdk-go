package commonnetbackup

import (
	"fmt"
	"time"

	clientnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/netbackup"
)

type JobStatus string

const (
	JobStatusCompleted JobStatus = "Completed"
	JobStatusFailed    JobStatus = "Failed"
	JobStatusRunning   JobStatus = "Running"
	JobStatusPending   JobStatus = "Pending"
)

// JobAPIResponse is the response structure for the Job API.
type JobAPIResponse struct {
	Data struct {
		ID     int    `json:"Id,omitempty"`
		Status string `json:"Status,omitempty"`
	} `json:"data,omitempty"`
}

// Refresh - Refreshes the job status.
func (j *JobAPIResponse) Refresh() error {
	c, err := clientnetbackup.New()
	if err != nil {
		return err
	}

	r, err := c.R().
		SetResult(&JobAPIResponse{}).
		SetError(&APIError{}).
		SetPathParams(map[string]string{
			"JobID": fmt.Sprintf("%d", j.Data.ID),
		}).
		Get("/v6/activities/{JobID}")
	if err != nil {
		return err
	}

	if r.IsError() {
		return ToError(r.Error().(*APIError))
	}

	*j = *r.Result().(*JobAPIResponse)

	return nil
}

// IsDone - Returns true if the job is done with a success.
func (j *JobAPIResponse) IsDone() bool {
	return j.Data.Status == string(JobStatusCompleted)
}

// Wait - Waits for the job to be done
// refreshInterval - The interval in seconds between each refresh
// timeout - The timeout in seconds.
func (j *JobAPIResponse) Wait(refreshInterval, timeout int) error {
	err := j.Refresh()
	if err != nil {
		return err
	}

	if j.IsDone() {
		return nil
	}

	ticker := time.NewTicker(time.Duration(refreshInterval) * time.Second)
	defer ticker.Stop()

	timeoutChan := time.After(time.Duration(timeout) * time.Second)

	for {
		select {
		case <-timeoutChan:
			return fmt.Errorf("timeout after %d seconds", timeout)
		case <-ticker.C:
			err := j.Refresh()
			if err != nil {
				return err
			}

			if j.IsDone() {
				return nil
			}
		}
	}
}
