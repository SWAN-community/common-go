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
	"encoding"
	"encoding/binary"
	"fmt"
	"time"
)

// ReadMarshaller reads the content into the unmarshaler instance.
func ReadMarshaller(b *bytes.Buffer, m encoding.BinaryUnmarshaler) error {
	v, err := ReadByteArray(b)
	if err != nil {
		return err
	}
	return m.UnmarshalBinary(v)
}

// WriteMarshaller writes the result of marshal binary call to the buffer.
func WriteMarshaller(b *bytes.Buffer, m encoding.BinaryMarshaler) error {
	v, err := m.MarshalBinary()
	if err != nil {
		return err
	}
	return WriteByteArray(b, v)
}

// ReadString reads a null (zero) terminated string from the byte buffer.
func ReadString(b *bytes.Buffer) (string, error) {
	s, err := b.ReadBytes(0)
	if err == nil {
		return string(s[0 : len(s)-1]), err
	}
	return "", err
}

// WriteString writes a null (zero) terminated string to the byte buffer.
func WriteString(b *bytes.Buffer, s string) error {
	l, err := b.WriteString(s)
	if err == nil {

		// Validate the number of bytes written matches the number of bytes in
		// the string.
		if l != len(s) {
			return fmt.Errorf(
				"mismatched lengths '%d' and '%d'",
				l,
				len(s))
		}

		// Write the null terminator.
		b.WriteByte(0)
	}
	return err
}

// ReadByteArray reads the first 4 bytes as an unsigned 32 bit integer to
// determine the length of the byte array contained in the following bytes.
func ReadByteArray(b *bytes.Buffer) ([]byte, error) {
	l, err := ReadUint32(b)
	if err != nil {
		return nil, err
	}
	return b.Next(int(l)), err
}

// ReadByteArray writes the length of the byte array as an unsigned 32 bit
// integer followed by the bytes.
func WriteByteArray(b *bytes.Buffer, v []byte) error {
	err := WriteUint32(b, uint32(len(v)))
	if err != nil {
		return err
	}
	return WriteByteArrayNoLength(b, v)
}

// ReadByteArrayNoLength reads the number of bytes specified into a new byte
// array.
func ReadByteArrayNoLength(b *bytes.Buffer, l int) ([]byte, error) {
	v := b.Next(l)
	if len(v) != l {
		return nil, fmt.Errorf("read '%d' bytes but expected '%d'", len(v), l)
	}
	return v, nil
}

// WriteByteArrayNoLength writes the byte array to the buffer without recording
// the length. Used with fixed length data.
func WriteByteArrayNoLength(b *bytes.Buffer, v []byte) error {
	l, err := b.Write(v)
	if err == nil {
		if l != len(v) {
			return fmt.Errorf(
				"mismatched lengths '%d' and '%d'",
				l,
				len(v))
		}
	}
	return err
}

// GetDateInMinutes returns the number of minutes that have elapsed since the
// IoDateBase epoch.
func GetDateInMinutes(t time.Time) uint32 {
	return uint32(t.Sub(IoDateBase).Minutes())
}

// GetTimeFromMinutes returns the date time from the minutes provided.
func GetDateFromMinutes(t uint32) time.Time {
	return IoDateBase.Add(time.Minute * time.Duration(t))
}

// ReadDateFromUInt32 reads the date from the buffer where the date is stored
// as the number of minutes as an unsigned 32 bit integer that have elapsed
// since the IoDateBase epoch.
func ReadDateFromUInt32(b *bytes.Buffer) (time.Time, error) {
	i, err := ReadUint32(b)
	if err != nil {
		return time.Time{}, err
	}
	return IoDateBase.Add(time.Duration(i) * time.Minute), nil
}

// WriteDateToUInt32 writes the date to the buffer as an unsigned 32 bit
// representing the number of minutes that have elapsed since the IoDateBase
// epoch.
func WriteDateToUInt32(b *bytes.Buffer, t time.Time) error {
	return WriteUint32(b, GetDateInMinutes(t))
}

// ReadByte reads the next byte from the buffer.
func ReadByte(b *bytes.Buffer) (byte, error) {
	d := b.Next(1)
	if len(d) != 1 {
		return 0, fmt.Errorf("'%d' bytes incorrect for byte", len(d))
	}
	return d[0], nil
}

// WriteByte writes the byte provided to the buffer.
func WriteByte(b *bytes.Buffer, i byte) error {
	return b.WriteByte(i)
}

// Read the next byte as a bool and returns the value.
func ReadBool(b *bytes.Buffer) (bool, error) {
	d := b.Next(1)
	if len(d) != 1 {
		return false, fmt.Errorf("'%d' bytes incorrect for bool", len(d))
	}
	return d[0] != 0, nil
}

// WriteBool writes the boolean value to the buffer as a byte.
func WriteBool(b *bytes.Buffer, v bool) error {
	var d byte
	if v {
		d = 1
	}
	return b.WriteByte(d)
}

// ReadUint32 reads an unsigned 32 bit integer from the buffer. The integer is
// stored in little endian format.
func ReadUint32(b *bytes.Buffer) (uint32, error) {
	d := b.Next(4)
	if len(d) != 4 {
		return 0, fmt.Errorf("'%d' bytes incorrect for Uint32", len(d))
	}
	return binary.LittleEndian.Uint32(d), nil
}

// ReadUint16 reads an unsigned 16 bit integer from the buffer. The integer is
// stored in little endian format.
func ReadUint16(b *bytes.Buffer) (uint16, error) {
	d := b.Next(2)
	if len(d) != 2 {
		return 0, fmt.Errorf("'%d' bytes incorrect for Uint16", len(d))
	}
	return binary.LittleEndian.Uint16(d), nil
}

// WriteUint16 reads an unsigned 16 bit integer from the buffer where the
// integer is stored in little endian format.
func WriteUint16(b *bytes.Buffer, i uint16) error {
	v := make([]byte, 2)
	binary.LittleEndian.PutUint16(v, i)
	l, err := b.Write(v)
	if err == nil {
		if l != len(v) {
			return fmt.Errorf(
				"Mismatched lengths '%d' and '%d'",
				l,
				len(v))
		}
	}
	return err
}

// WriteUint32 reads an unsigned 32 bit integer from the buffer where the
// integer is stored in little endian format.
func WriteUint32(b *bytes.Buffer, i uint32) error {
	v := make([]byte, 4)
	binary.LittleEndian.PutUint32(v, i)
	l, err := b.Write(v)
	if err == nil {
		if l != len(v) {
			return fmt.Errorf(
				"mismatched lengths '%d' and '%d'",
				l,
				len(v))
		}
	}
	return err
}
