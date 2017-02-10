// Package httpmessageconverter provides a strategy interface
// that specifies a converter that converts to and from
// HTTP requests and responses.
// Copyright 2016 Don B. Stringham. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.
//
// @author Don B. Stringham <donbstringham@gmail.com>
//
package msgconv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanMediaType_Fail(t *testing.T) {
	// arrange
	c := "Content"
	// act
	v := CleanMediaType(c)
	// assert
	assert.Equal(t, "Content", v)
}

func TestCleanMediaType_Success0(t *testing.T) {
	// arrange
	c := "Content-Type:application/json"
	// act
	v := CleanMediaType(c)
	// assert
	assert.Equal(t, "Content-Type:application/json", v)
}

func TestCleanMediaType_Success1(t *testing.T) {
	// arrange
	c := "Accept:application/json;Content-Type:application/json"
	// act
	v := CleanMediaType(c)
	// assert
	assert.Equal(t, "Accept:application/json", v)
}

func TestMsgConv(t *testing.T) {
	assert.Equal(t, "ISO-8859-1", CharsetText[CharsetIso88591])
	assert.Equal(t, "US-ASCII", CharsetText[CharsetUsAscii])
	assert.Equal(t, "UTF-8", CharsetText[CharsetUTF8])
}
