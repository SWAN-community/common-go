/* ****************************************************************************
 * Copyright 2022 51 Degrees Mobile Experts Limited (51degrees.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 * ***************************************************************************/
package common

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// HTTPTest returns a test response after having processed the handler, method,
// URL, and body provided.
// t testing instance
// method HTTP method
// host
// url HTTP url
// values query values
// body data
// handler HTTP handler being tested
func HTTPTest(
	t *testing.T,
	method string,
	host string,
	url string,
	values url.Values,
	handler func(w http.ResponseWriter, r *http.Request)) *httptest.ResponseRecorder {

	// Add the values to the body if not a GET.
	var body io.Reader
	if method != "GET" && values != nil {
		body = strings.NewReader(values.Encode())
	} else {
		body = nil
	}

	// Get the new request.
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}

	// Set the host and query string parameters if a GET.
	req.Host = host
	if method == "GET" && values != nil {
		req.URL.RawQuery = values.Encode()
	}

	// Call the handler and return the response.
	rr := httptest.NewRecorder()
	http.HandlerFunc(handler).ServeHTTP(rr, req)
	return rr
}

func TestCompareDate(t *testing.T, a time.Time, b time.Time) {
	if a.Year() != b.Year() {
		fmt.Printf("Year %d != %d", a.Year(), b.Year())
		t.Fail()
	}
	if a.Month() != b.Month() {
		fmt.Printf("Month %d != %d", a.Month(), b.Month())
		t.Fail()
	}
	if a.Day() != b.Day() {
		fmt.Printf("Day %d != %d", a.Day(), b.Day())
		t.Fail()
	}
}
