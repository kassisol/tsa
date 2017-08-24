package httputils

import (
	"net/url"
)

func QueryParams2Filters(val url.Values) map[string]string {
	result := make(map[string]string)

	for k, v := range val {
		result[k] = v[0]
	}

	return result
}
