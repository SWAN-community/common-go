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
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestReturnServerError simulates a server error.
func TestReturnServerError(t *testing.T) {
	err := errors.New("A")
	t.Run("with error", func(t *testing.T) {
		testReturnServerError(t, err)
	})
	t.Run("without error", func(t *testing.T) {
		testReturnServerError(t, nil)
	})
}

// TestReturnApplicationError simulates common application errors where a
// message and status code is returned in the HTTP response and an associated
// error could be logged.
func TestReturnApplicationError(t *testing.T) {
	message := "A"
	err := errors.New("B")
	t.Run("log all", func(t *testing.T) {
		testReturnApplicationError(
			t, true, message, err, http.StatusBadRequest)
	})
	t.Run("no log all", func(t *testing.T) {
		testReturnApplicationError(
			t, false, message, err, http.StatusBadRequest)
	})
	t.Run("log no error", func(t *testing.T) {
		testReturnApplicationError(
			t, true, message, nil, http.StatusBadRequest)
	})
	t.Run("no log no error", func(t *testing.T) {
		testReturnApplicationError(
			t, false, message, nil, http.StatusBadRequest)
	})
}

func testReturnServerError(t *testing.T, err error) {
	rr := HTTPTest(t, "GET", "/test", nil, func(w http.ResponseWriter, r *http.Request) {
		ReturnServerError(w, err)
	})
	validateCode(t, rr, http.StatusInternalServerError)
	validateMessage(t, rr, serverErrorMessage)
}

func testReturnApplicationError(
	t *testing.T,
	log bool,
	message string,
	err error,
	code int) {
	rr := HTTPTest(t, "GET", "/test", nil, func(w http.ResponseWriter, r *http.Request) {
		ReturnApplicationError(w, &HttpError{
			Request: r,
			Log:     log,
			Message: message,
			Code:    code,
			Error:   err})
	})
	validateCode(t, rr, code)
	validateMessage(t, rr, message)
}

func validateMessage(t *testing.T, rr *httptest.ResponseRecorder, message string) {
	if strings.Contains(rr.Body.String(), message) == false {
		t.Errorf("handler returned unexpected body: got '%v' expected '%v'",
			rr.Body.String(), message)
	}
}

func validateCode(t *testing.T, rr *httptest.ResponseRecorder, code int) {
	if rr.Code != code {
		t.Errorf(
			"handler returned wrong status code: got '%v' expected '%v'",
			rr.Code, code)
	}
}
