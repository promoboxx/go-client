package service

import (
	"net/http"
	"strings"
)

type Sort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

func (s Sort) Valid(fieldWhiteList []string) bool {
	if s.Field == "" || !contains(fieldWhiteList, s.Field) {
		return false
	}

	switch s.Direction {
	case "asc",
		"desc":
	// no-op
	default:
		return false
	}

	return true
}

// PagingParams represents paging and sorting parameter values
type PagingParams struct {
	PageSize   *int32 `json:"page_size"`
	Offset     *int32 `json:"offset"`
	SortFields []Sort `json:"sort_fields"`
}

// ParsePagingParams retrieves paging params from the request, allows for whitelisting sort field values
func ParsePagingParams(r *http.Request, defaults PagingParams, sortFieldsWhitelist []string) (PagingParams, error) {
	paging := defaults

	// offset (page number)
	offset, err := Int32PointerFromQueryParam(r, "offset")
	if err != nil {
		return paging, err
	}
	if offset != nil {
		paging.Offset = offset
	}

	// page size
	pageSize, err := Int32PointerFromQueryParam(r, "page_size")
	if err != nil {
		return paging, err
	}
	if pageSize != nil {
		paging.PageSize = pageSize
	}

	// Sort fields
	sorts := strings.Split(r.URL.Query().Get("sort"), ",")
	for _, sort := range sorts {
		if sort != "" {
			d := strings.Split(sort, ":")
			s := Sort{
				Field:     strings.ToLower(d[0]),
				Direction: strings.ToLower(d[1]),
			}

			if valid := s.Valid(sortFieldsWhitelist); valid {
				paging.SortFields = append(paging.SortFields, s)
			}
		}
	}

	return paging, nil
}

// Contains tells whether a contains x.
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
