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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// HTTPTest returns a test response after having processed the handler, method, URL, and body provided.
// t testing instance
// method HTTP method
// url HTTP url
// body HTTP query string or body data
// handler HTTP handler being tested
func HTTPTest(
	t *testing.T,
	method string,
	url string,
	body io.Reader,
	handler func(w http.ResponseWriter, r *http.Request)) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(handler)
	httpHandler.ServeHTTP(rr, req)
	return rr
}
