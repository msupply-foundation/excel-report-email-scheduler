package server

import "fmt"

type RequestBodyError struct {
	m string
}

func (e *RequestBodyError) Error() string {
	return e.m
}

func NewRequestBodyError(err error, bodyShouldBe string) *RequestBodyError {
	m := fmt.Sprintf("Error: %s\nExpecting a body with the shape: %s", err.Error(), bodyShouldBe)
	return &RequestBodyError{m: m}
}
