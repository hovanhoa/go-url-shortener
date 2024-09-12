package url

import "net/url"

func IsValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return false
	}

	return true
}
