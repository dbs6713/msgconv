// Copyright 2017 Intellectual Reserve, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.
//
// @author Don B. Stringham <stringhamdb@ldschurch.org>
//
package msgconv

import (
	"testing"
	"reflect"

	"github.com/stretchr/testify/assert"
)

func TestNewPBHTTPConverter(t *testing.T) {
	// arrange and act
	v := NewPBHTTPConverter()
	// assert
	assert.Equal(t, "*msgconv.PBHTTPConverter", reflect.TypeOf(v).String())
}
