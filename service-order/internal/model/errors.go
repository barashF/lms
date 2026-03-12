package model

import "errors"

var (
	ErrNoOrderFound            = errors.New("no order found")
	ErrMessageAlreadyProcessed = errors.New("message already processed")
)
