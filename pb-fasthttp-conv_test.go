// Package msgconv provides a strategy interface
// that specifies a converter that converts to and from
// fasthttp HTTP requests and responses.
// Copyright 2016 Don B. Stringham. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.
//
// @author Don B. Stringham <donbstringham@gmail.com>
//
package msgconv

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPBFasthttpConverter(t *testing.T) {
	// arrange and act
	v := NewPBFasthttpConverter()
	// assert
	assert.Equal(t, "*msgconv.PBFasthttpConverter", reflect.TypeOf(v).String())
}

func TestPBFasthttpConverter_CanRead_Fail(t *testing.T) {
	// arrange
	c := NewPBFasthttpConverter()
	// act and assert
	assert.False(t, c.CanRead(-1))
}

func TestPBFasthttpConverter_CanRead_Success(t *testing.T) {
	// arrange
	c := NewPBFasthttpConverter()
	// act and assert
	assert.True(t, c.CanRead(MediatypeJSON))
}

func TestPBFasthttpConverter_CanWrite(t *testing.T) {
	t.Skip()
}

func TestPBFasthttpConverter_CleanMediaType(t *testing.T) {
	t.Skip()
}

func TestPBFasthttpConverter_GetMediaType(t *testing.T) {
	t.Skip()
}

func TestPBFasthttpConverter_ReadRequest(t *testing.T) {
	t.Skip()
}

func TestPBFasthttpConverter_ReadRequestCtx(t *testing.T) {
	t.Skip()
}

func TestPBFasthttpConverter_ReadResponse(t *testing.T) {
	t.Skip()
}

func TestPBFasthttpConverter_WriteRequest(t *testing.T) {
	t.Skip()
}

func TestPBFasthttpConverter_WriteRequestCtx(t *testing.T) {
	t.Skip()
}

func TestPBFasthttpConverter_WriteResponse(t *testing.T) {
	t.Skip()
}

func TestPBFasthttpConverter_WriteResponseCtx(t *testing.T) {
	t.Skip()
}
