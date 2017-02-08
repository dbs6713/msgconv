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

func TestMsgConv(t *testing.T) {
	assert.Equal(t, "ISO-8859-1", CharsetText[CharsetIso88591])
	assert.Equal(t, "US-ASCII", CharsetText[CharsetUsAscii])
	assert.Equal(t, "UTF-8", CharsetText[CharsetUTF8])
}
