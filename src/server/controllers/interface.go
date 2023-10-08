package controllers

type RequestValidator interface {
	ValidateRequest() error
}
