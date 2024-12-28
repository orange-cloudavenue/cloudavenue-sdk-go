package commoncloudavenue

import (
	"context"
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
		return j.Wait(refreshInterval, 90)
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
			return fmt.Errorf("timeout reached")
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
