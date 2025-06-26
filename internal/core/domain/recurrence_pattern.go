package domain

type RecurrencePattern struct {
	ID          string `json:"id"`
	IngressID    string `json:"ingress_id"`
	Frequency     Frequency `json:"frequency"`
	Interval     int `json:"interval"`
	Amount        float64 `json:"amount"`
	EndDate      *string `json:"endDate,omitempty"`
}

type Frequency string

const (
    Daily Frequency = "daily"
    Weekly Frequency = "weekly"
    Monthly Frequency = "monthly"
	Yearly Frequency = "yearly"
)