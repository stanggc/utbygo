package main

// #include "preamble.h"
import "C"

import "crypto/tls"
import "net/http"
import url_ "net/url"
import "bytes"
import "fmt"
import "io"
import "unsafe"

//export Fetch
func Fetch(method, urlC, headers, content *C.char, contentLength C.size_t, options C.struct_FetchOptions, outResp *C.struct_Response) C.enum_Status {
	url, err := url_.Parse(C.GoString(urlC))
	if err != nil {
		return C.E_URL_PARSE
	}

	skipVerifyTLSCert := bool(options.SkipVerifyTLSCert)
	client := &http.Client{}
	if skipVerifyTLSCert {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: skipVerifyTLSCert,
			},
		}
		client.Transport = transport
	}
	req := &http.Request{
		Method: C.GoString(method),
		URL:    url,
	}

	// Set up headers.
	if headers != nil {
		req.Header = StringToHeaders(C.GoString(headers), nil)
	}

	// Set up content, if provided.
	if content != nil {
		req.Header.Add("content-length", fmt.Sprintf("%d", contentLength))
		buf := io.NopCloser(bytes.NewBuffer(C.GoBytes(unsafe.Pointer(content), C.int(contentLength))))
		req.Body = buf
	}

	resp, err := client.Do(req)
	if err != nil {
		return C.E_FETCH
	}
	respContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return C.E_READ
	}
	// Convert response header into a string.
	respHdrs := HeadersToString(resp.Header)

	if outResp != nil {
		*outResp = C.struct_Response{
			StatusCode:    C.int(resp.StatusCode),
			Headers:       C.CString(respHdrs),
			ContentLength: C.size_t(len(respContent)),
			Content:       (*C.char)(C.CBytes(respContent)),
		}
	}
	return C.E_OK
}
