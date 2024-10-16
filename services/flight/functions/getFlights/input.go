package main

type (
	Input struct {
		Name string `json:"name" validate:"required"`
	}
)
