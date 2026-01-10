package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

func getQueryInt(r *http.Request, name string, defaultValue int) (int, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return defaultValue, nil
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid '%s' parameter", name)
	}

	return parsedValue, nil
}
