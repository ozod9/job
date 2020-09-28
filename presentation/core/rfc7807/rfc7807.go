package rfc7807

import (
	"encoding/json"
	"net/http"
)

type Problem struct {
	Type     string  `json:"type"`
	Title    *string `json:"title"`
	Status   int     `json:"status"`
	Detail   string  `json:"detail"`
	Instance *string `json:"instance"`
	Errors   []Error `json:"errors"`
}

type Error struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

func NewProblem() *Problem {
	problem := new(Problem)
	problem.Errors = make([]Error, 0)
	return problem
}

func (problem *Problem) SetType(errorType string) *Problem {
	problem.Type = errorType
	return problem
}

func (problem *Problem) AppendError(name, reason string) *Problem {
	err := Error{}
	err.Name = name
	err.Reason = reason
	problem.Errors = append(problem.Errors, err)
	return problem
}

func (problem *Problem) SetStatus(status int) *Problem {
	problem.Status = status
	return problem
}

func (problem *Problem) Write(w http.ResponseWriter) error {
	body, err := json.Marshal(problem)
	if err != nil {
		return err
	}
	w.WriteHeader(problem.Status)
	w.Header().Set("Content-Type", "application/problem+json")
	if _, err := w.Write(body); err != nil {
		return err
	}
	return nil
}
