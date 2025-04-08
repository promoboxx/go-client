package glitch

import "fmt"

type GQLProblem struct {
	PublicMsg string `json:"message,omitempty"`
	ErrorCode string `json:"code,omitempty"`
}

func (g *GQLProblem) Error() string {
	return fmt.Sprintf("Error Code [%s] - %s", g.ErrorCode, g.PublicMsg)
}
