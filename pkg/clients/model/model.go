package model

type (
	ClientOpts interface {
		Validate() error
	}
)
