// Package msgconv provides a strategy interface
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
	"github.com/golang/protobuf/proto"
	"strings"
)

const (
	CharsetIso88591 = iota
	CharsetUsAscii
	CharsetUTF8
	MediatypeJSON = iota
	MediatypePB
	MessageNative = iota
	MessageFasthttp
	Accept        = "Accept"
	CacheControl  = "Cache-Control"
	ContentLength = "Content-Length"
	ContentType   = "Content-Type"
)

var (
	MediatypeText = map[int]string{
		MediatypeJSON: "application/json",
		MediatypePB:   "application/x-protobuf",
	}
	CharsetText = map[int]string{
		CharsetIso88591: "ISO-8859-1",
		CharsetUsAscii:  "US-ASCII",
		CharsetUTF8:     "UTF-8",
	}
)

type MessageConverter interface {
	CanRead(mediatype int) bool
	CanWrite(mediatype int) bool
	GetMediaType(mediatype string) (int, error)
	ReadRequest(req *interface{}, v interface{}) error
	ReadResponse(resp *interface{}, v interface{}) error
	WriteRequest(method string, URI string, contentType int, msg proto.Message) (interface{}, error)
	WriteResponse(v interface{}) (interface{}, error)
}

func CleanMediaType(mediatype string) string {
	lastIdx := strings.LastIndex(mediatype, ";")
	if lastIdx != -1 {
		return mediatype[0:lastIdx]
	}
	return mediatype
}
