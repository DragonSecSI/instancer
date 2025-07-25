package api

import (
	"fmt"
	"net/http"
	"strconv"
)

type ApiPagination struct {
	GetPagination func(r *http.Request) (int, int, error)
}

func getPagination(r *http.Request) (int, int, error) {
	page := 1
	pagesize := 10

	if r.URL.Query().Get("page") != "" {
		var err error
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			return 0, 0, fmt.Errorf("invalid page number")
		}
	}

	if r.URL.Query().Get("pagesize") != "" {
		var err error
		pagesize, err = strconv.Atoi(r.URL.Query().Get("pagesize"))
		if err != nil || pagesize < 1 {
			return 0, 0, fmt.Errorf("invalid pagesize number")
		}
	}

	return page, pagesize, nil
}
