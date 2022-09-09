/* ****************************************************************************
 * Copyright 2020 51 Degrees Mobile Experts Limited (51degrees.com)
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
	"bytes"
	"fmt"
	"testing"
	"time"
)

// TestIoTime verifies the storage of time as an unsigned 32 bit integer.
func TestIoTime(t *testing.T) {
	d := time.Now().UTC()
	var b bytes.Buffer
	err := WriteDateToUInt32(&b, d)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	i := b.Bytes()
	c := bytes.NewBuffer(i)
	r, err := ReadDateFromUInt32(c)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	TestCompareDate(t, r, d)
}
