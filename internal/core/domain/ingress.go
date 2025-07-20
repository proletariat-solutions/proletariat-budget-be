package domain

import "errors"

// Ingress domain errors
var (
	ErrIngressNotFound           = errors.New("ingress not found")
	ErrRecurrencePatternNotFound = errors.New("recurrence pattern not found")
)
