package main

type webError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
