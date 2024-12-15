package request

import (
	"net/http"
)

type WithAuthorize interface {
	Authorize(r *http.Request) error
}

type WithPrepare interface {
	Prepare(r *http.Request) error
}

type WithRules interface {
	Rules(r *http.Request) map[string]string
}

type WithFilters interface {
	Filters(r *http.Request) map[string]string
}

type WithMessages interface {
	Messages(r *http.Request) map[string]string
}
