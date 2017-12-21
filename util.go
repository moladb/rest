package rest

import (
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

func ginHandlerFunc(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func trimDuplicateRune(s string, r rune) string {
	sLen := len(s)
	if sLen == 0 {
		return s
	}

	result := make([]byte, sLen)
	inputI := 0
	preI := 0
	preR, width := utf8.DecodeRuneInString(s)
	if preR == utf8.RuneError {
		panic(fmt.Sprintf("string contains invalid utf-8 code point\n{% x}\n", s))
	}
	curI := preI + width
	copy(result, s[preI:curI])
	inputI += width

	for curI < sLen {
		curR, width := utf8.DecodeRuneInString(s[curI:])
		if curR == utf8.RuneError {
			panic(fmt.Sprintf("string contains invalid utf-8 code point\n{% x}\n", s[curI:]))
		}
		if preR == r && curR == r {
			curI += width
		} else {
			preR = curR
			preI = curI
			curI += width
			copy(result[inputI:], s[preI:curI])
			inputI += width
		}
	}

	return string(result[:inputI])
}

func fullResourcePath(group, resource string) string {
	return normalizePath(fmt.Sprintf("%s/%s", group, resource))
}

func normalizePath(s string) string {
	// assume s is all assic
	return trimDuplicateRune("/"+strings.TrimLeft(s, "/"), '/')
}
