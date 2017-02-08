// Package core provides the base functions
// Copyright 2016 Don B. Stringham. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.
//
// @author Don B. Stringham <donbstringham@gmail.com>
//
package core

import "testing"

func TestVersion(t *testing.T) {
	equals(t, Version, "")
	equals(t, Name, "")
	equals(t, Build, "")
	equals(t, BuildTime, "")
}
