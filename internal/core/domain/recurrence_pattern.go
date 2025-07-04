package domain

import "ghorkov32/proletariat-budget-be/openapi"

type RecurrencePattern struct {
	ID        string    `json:"id"`
	IngressID string    `json:"ingress_id"`
	Frequency Frequency `json:"frequency"`
	Interval  int       `json:"interval"`
	Amount    float32   `json:"amount"`
	EndDate   *string   `json:"endDate,omitempty"`
}

type Frequency string

const (
	Daily   Frequency = "daily"
	Weekly  Frequency = "weekly"
	Monthly Frequency = "monthly"
	Yearly  Frequency = "yearly"
)

func getFrequencyFromOAPI(frequency openapi.RecurrencePatternFrequency) Frequency {
	switch frequency {
	case openapi.RecurrencePatternFrequency(openapi.Daily):
		return Daily
	case openapi.RecurrencePatternFrequency(openapi.Weekly):
		return Weekly
	case openapi.RecurrencePatternFrequency(openapi.Monthly):
		return Monthly
	case openapi.RecurrencePatternFrequency(openapi.Yearly):
		return Yearly
	default:
		return Daily
	}
}

func NewRecurrencePattern(recurrencePattern *openapi.RecurrencePattern, ingressID string) *RecurrencePattern {
	endDate := (*recurrencePattern.EndDate).String()
	return &RecurrencePattern{
		ID:        *recurrencePattern.Id,
		IngressID: ingressID,
		Frequency: getFrequencyFromOAPI(*recurrencePattern.Frequency),
		Interval:  *recurrencePattern.Interval,
		Amount:    *recurrencePattern.Amount,
		EndDate:   &endDate,
	}
}
