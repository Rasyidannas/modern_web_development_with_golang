package models

type Result struct {
	ID string `json:"id"`
	RunnerID string `json:"runner_id"`
	RaceResult string 
	Location string `json:"location"`
	Position int `json:"position,omitempty"`
	Year int `json:"year"`
}
