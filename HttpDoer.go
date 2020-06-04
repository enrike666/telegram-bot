package main

import "net/http"

type HttpDoer interface {
	Do(*http.Request) (*http.Response, error)
}
