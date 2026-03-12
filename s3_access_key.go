// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s3utils

type AccessKey string

func (a AccessKey) String() string {
	return string(a)
}
