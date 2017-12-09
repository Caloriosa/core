package tools

import (
	"net/url"
	"strconv"
)

const DEFAULT_LIMIT = 50

func GetPagination(values url.Values) (int, int, error) {
	var page, limit int
	var err error

	if val, ok := values["page"]; ok  {
		page, err = strconv.Atoi(val[0])
		if err != nil {
			return 0, 0, err
		}
	} else {
		page = 0
	}

	if val, ok := values["limit"]; ok {
		limit, err = strconv.Atoi(val[0])
		if err != nil {
			return 0, 0, nil
		}
	} else {
		limit = DEFAULT_LIMIT
	}

	return page, limit, nil
}
