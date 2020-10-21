package glitch

import "fmt"

// HTTPProblem should be used as the response in case of an error during an HTTP request.
// It implements the https://datatracker.ietf.org/doc/rfc7807 spec with these additional fields:
// 		code: meant to be machine readable and give clients enough information to handle the error appropriately
//		is_transient: meant to inform clients that the problem is considered transient and could be retried
// swagger:model HTTPProblem
type HTTPProblem struct {
	Type        string `json:"type,omitempty"`
	Title       string `json:"title,omitempty"`
	Status      int    `json:"status,omitempty"`
	Detail      string `json:"detail,omitempty"`
	Instance    string `json:"instance,omitempty"`
	Code        string `json:"code,omitempty"`
	IsTransient bool   `json:"is_transient"`
}

func (h HTTPProblem) Error() string {
	transient := ""
	if h.IsTransient {
		transient = " (transient)"
	}
	return fmt.Sprintf("HTTPProblem: [%d - %s%s] - %s - %s", h.Status, h.Code, transient, h.Title, h.Detail)
}
