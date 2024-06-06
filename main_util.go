package main

import "C"
import "fmt"
import "net/http"
import "strings"
import "time"

func HeadersToString(hdrs http.Header) string {
	if hdrs == nil {
		return ""
	}

	hdrsStr := ""
	for k, v := range hdrs {
		for _, vv := range v {
			hdrsStr = fmt.Sprintf("%s%s: %s\r\n", hdrsStr, k, vv)
		}
	}

	return hdrsStr
}

func StringToHeaders(hdrsStr string, outHeaders http.Header) http.Header {
	lines := strings.Split(hdrsStr, "\r\n")
	var out http.Header
	if outHeaders != nil {
		out = outHeaders
	} else {
		out = make(http.Header)
	}

	for _, line := range lines {
		vals := strings.Split(line, ": ")
		if len(vals) < 2 {
			if vals[0] != "" {
				out.Add(vals[0], "")
			}
		} else {
			out.Add(vals[0], vals[1])
		}
	}

	return out
}

//export Sleep
func Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
