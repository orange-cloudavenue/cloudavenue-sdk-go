package commoncloudavenue

import (
	"fmt"
	"time"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
)

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

// JobCreatedAPIResponse - This is the response structure for the JobCreatedAPIResponse
type JobCreatedAPIResponse struct {
	Message string `json:"message"`
	JobID   string `json:"jobId"`
}

// JobStatus - This is the response structure for the JobStatus
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

// GetJobStatus - Returns the status of a job
func (j *JobCreatedAPIResponse) GetJobStatus() (response *JobStatus, err error) {
	response.JobID = j.JobID
	if err := response.Refresh(); err != nil {
		return nil, err
	}

	return response, response.Refresh()
}

// Refresh - Refreshes the job status
func (j *JobStatus) Refresh() error {
	jobID := j.JobID

	c, err := clientcloudavenue.New()
	if err != nil {
		return err
	}

	r, err := c.R().
		SetResult(&[]JobStatus{}).
		SetError(&APIErrorResponse{}).
		SetPathParams(map[string]string{
			"JobID": j.JobID,
		}).
		Get("/api/customers/v1.0/jobs/{JobID}")
	if err != nil {
		return err
	}

	if r.IsError() {
		return ToError(r.Error().(*APIErrorResponse))
	}

	x := *r.Result().(*[]JobStatus)
	*j = x[0]

	j.JobID = jobID

	return nil
}

// IsDone - Returns true if the job is done with a success
func (j *JobStatus) IsDone() bool {
	return j.Status == DONE
}

// Wait - Waits for the job to be done
// refreshInterval - The interval in seconds between each refresh
// timeout - The timeout in seconds
func (j *JobStatus) Wait(refreshInterval, timeout int) error {
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
