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
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Message to return in the HTTP response when a server error occurs.
const serverErrorMessage = "Internal server error"

// HttpError associated with HTTP handlers.
type HttpError struct {
	Request *http.Request // the HTTP request
	Log     bool          // true if the error should be written to the log
	Message string        // message to return in the HTTP response
	Code    int           // HTTP status code for the response
	Error   error         // the server error - never sent in the response
}

// ReturnApplicationError handles HTTP application errors consistently.
// writer for the response
// err details of the error
func ReturnApplicationError(writer http.ResponseWriter, err *HttpError) {
	ReturnError(writer, err)
}

// ReturnServerError handles HTTP server errors consistently ensuring they are
// output to the logger.
// writer for the response
// message to be sent in the response
// err the error to be logged and included in the response if debug is true
func ReturnServerError(writer http.ResponseWriter, err error) {
	ReturnError(writer, &HttpError{
		Log:     true,
		Message: serverErrorMessage,
		Code:    http.StatusInternalServerError,
		Error:   err})
}

// ReturnError handles all HTTP errors consistently.
// writer for the response
// err details of the error
func ReturnError(writer http.ResponseWriter, err *HttpError) {
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(writer, err.Message, err.Code)
	err.logError()
}

// logError if the log flag is set to true using a format to make it easier
// for operators to understand the cause of the error.
func (err *HttpError) logError() {
	if err.Log {
		var b strings.Builder
		b.WriteString("HTTP Error\r\n")
		b.WriteString("\tMessage: " + err.Message + "\r\n")
		b.WriteString("\tCode   : " + strconv.Itoa(err.Code) + "\r\n")
		if err.Error != nil {
			b.WriteString("\tError  : " + err.Error.Error() + "\r\n")
		}
		if err.Request != nil {
			b.WriteString("\tMethod : " + err.Request.Method + "\r\n")
			b.WriteString("\tURL    : " + err.Request.URL.String() + "\r\n")
		}
		log.Print(b.String())
	}
}
