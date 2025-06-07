package utils

import (
	"net/http"
	"strings"
)

func ParseCookieStr(cookieStr string) []*http.Cookie {
	var cookies []*http.Cookie
	if cookieStr == "" {
		return cookies
	}
	for _, partStr := range strings.Split(cookieStr, ";") {
		parts := strings.SplitN(partStr, "=", 2)
		if len(parts) != 2 {
			continue
		}
		cookies = append(cookies, &http.Cookie{
			Name:  strings.TrimSpace(parts[0]),
			Value: strings.TrimSpace(parts[1]),
		})
	}
	return cookies
}
