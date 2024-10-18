package main

type (
	FlightInput struct {
		Name string `json:"name" validate:"required"`
	}
	Input struct {
		Flight FlightInput `json:"flight" validate:"required"`
	}
)
