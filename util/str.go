// WUTONG, Application Management Platform
// Copyright (C) 2014-2017 Wutong Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Wutong,
// one or multiple Commercial Licenses authorized by Wutong Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package util

import (
	"reflect"
	"unsafe"
)

// ToString []byte to string
// BenchmarkToString-4     30000000                42.0 ns/op
func ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ToByte string to []byte
// BenchmarkToByte-4       2000000000               0.36 ns/op
func ToByte(v string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&v))
	return *(*[]byte)(unsafe.Pointer(sh))
}

// StringArrayContains string array contains
func StringArrayContains(list []string, source string) bool {
	for _, l := range list {
		if l == source {
			return true
		}
	}
	return false
}

// Reverse reverse sort string array
func Reverse(source []string) []string {
	for i := 0; i < len(source)/2; i++ {
		source[i], source[len(source)-i-1] = source[len(source)-i-1], source[i]
	}
	return source
}
