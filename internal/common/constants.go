package common

import (
	"time"
)

const (
	APPNAME = "tmpl-com-rest-api"

	RequestID = "request_id"

	// TODO: Update acondingly
	Domain = "my_api_domain"
	Create = "create"

	RetryWaitMin = 800 * time.Microsecond
	RetryWaitMax = 1200 * time.Microsecond
)
