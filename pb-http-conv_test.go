// Copyright 2017 Intellectual Reserve, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.
//
// @author Don B. Stringham <stringhamdb@ldschurch.org>
//
package msgconv

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/protobuf/jsonpb/jsonpb_test_proto"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

type respTest struct {
	Raw  string
	Resp http.Response
}

const (
	JSON = iota
	JSON_BAD
)

func dummyReq(method string) *http.Request {
	return &http.Request{Method: method, Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1}
}

var respTests = []respTest{
	{
		"HTTP/1.1 200 OK\r\nAccept:application/json; charset=UTF-8\r\nContent-Type:application/json; charset=UTF-8\r\n\r\n{\"o_string\":\"test\"}",
		http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Request:    dummyReq("GET"),
			Header: http.Header{
				Accept:      {MediatypeText[MediatypeJSON]},
				ContentType: {MediatypeText[MediatypeJSON]},
			},
			Close:         true,
			ContentLength: int64(19),
		},
	},
	{
		"HTTP/1.1 200 OK\r\nAccept:text/plain; charset=UTF-8\r\nContent-Type:text/plain; charset=UTF-8\r\n\r\n{\"o_string\":\"test\"}",
		http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Request:    dummyReq("GET"),
			Header: http.Header{
				Accept:      {MediatypeText[MediatypeJSON]},
				ContentType: {MediatypeText[MediatypeJSON]},
			},
			Close:         true,
			ContentLength: int64(19),
		},
	},
}

func TestNewPBHTTPConverter(t *testing.T) {
	// arrange and act
	v := NewPBHTTPConverter()
	// assert
	assert.Equal(t, "*msgconv.PBHTTPConverter", reflect.TypeOf(v).String())
}

func TestPBHTTPConverter_CanRead_Fail(t *testing.T) {
	// arrange
	c := NewPBHTTPConverter()
	// act
	assert.False(t, c.CanRead(-1))
}

func TestPBHTTPConverter_CanRead_Success(t *testing.T) {
	// arrange
	c := NewPBHTTPConverter()
	// act
	assert.True(t, c.CanRead(MediatypePB))
}

func TestPBHTTPConverter_CanWrite_Fail(t *testing.T) {
	// arrange
	c := NewPBHTTPConverter()
	// act
	assert.False(t, c.CanWrite(-1))
}

func TestPBHTTPConverter_CanWrite_Success(t *testing.T) {
	// arrange
	c := NewPBHTTPConverter()
	// act
	assert.True(t, c.CanWrite(MediatypePB))
}

func TestPBHTTPConverter_GetMediaType_Fail(t *testing.T) {
	// arrange
	c := NewPBHTTPConverter()
	// act
	_, err := c.GetMediaType("text/plain")
	// assert
	assert.NotNil(t, err)
	assert.Equal(t, "media type not supported", err.Error())
}

func TestPBHTTPConverter_GetMediaType_Success(t *testing.T) {
	// arrange
	c := NewPBHTTPConverter()
	// act
	v, err := c.GetMediaType("application/json")
	// assert
	assert.Nil(t, err)
	assert.Equal(t, MediatypeJSON, v)
}

func TestPBHTTPConverter_ReadRequest_Fail(t *testing.T) {
	// arrange
	var m map[string]interface{}
	p := jsonpb.Simple{OString: proto.String("test")}
	j, _ := json.Marshal(p)
	r, _ := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewBuffer(j))
	r.Header.Add(Accept, "text/plain; charset=UTF-8")
	c := NewPBHTTPConverter()
	// act
	err := c.ReadRequest(r, &m)
	// assert
	assert.NotNil(t, err)
	assert.Equal(t, "media type not supported", err.Error())
}

func TestPBHTTPConverter_ReadRequest_JSON(t *testing.T) {
	// arrange
	var m map[string]interface{}
	p := jsonpb.Simple{OString: proto.String("test")}
	j, _ := json.Marshal(p)
	r, _ := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewBuffer(j))
	r.Header.Add(Accept, MediatypeText[MediatypeJSON])
	c := NewPBHTTPConverter()
	// act
	err := c.ReadRequest(r, &m)
	// assert
	assert.Nil(t, err)
	assert.Equal(t, "test", m["o_string"])
}

func TestPBHTTPConverter_ReadRequest_PB(t *testing.T) {
	// arrange
	pm, _ := proto.Marshal(&jsonpb.Simple{OString: proto.String("test")})
	r, _ := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewBuffer(pm))
	r.Header.Add(Accept, MediatypeText[MediatypePB])
	c := NewPBHTTPConverter()
	m := jsonpb.Simple{}
	// act
	err := c.ReadRequest(r, &m)
	// assert
	assert.Nil(t, err)
	assert.Equal(t, "test", m.GetOString())
}

func TestPBHTTPConverter_ReadResponse_Fail(t *testing.T) {
	// arrange
	var m map[string]interface{}
	resp, _ := http.ReadResponse(bufio.NewReader(strings.NewReader(respTests[JSON_BAD].Raw)), respTests[JSON_BAD].Resp.Request)
	c := NewPBHTTPConverter()
	// act
	err := c.ReadResponse(resp, &m)
	// assert
	assert.NotNil(t, err)
	assert.Equal(t, "media type not supported", err.Error())
}

func TestPBHTTPConverter_ReadResponse_JSON(t *testing.T) {
	// arrange
	var m map[string]interface{}
	resp, _ := http.ReadResponse(bufio.NewReader(strings.NewReader(respTests[JSON].Raw)), respTests[JSON].Resp.Request)
	c := NewPBHTTPConverter()
	// act
	err := c.ReadResponse(resp, &m)
	// assert
	assert.Nil(t, err)
	assert.Equal(t, "test", m["o_string"])
}

func TestPBHTTPConverter_ReadResponse_PB(t *testing.T) {
	// arrange
	converter := func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		w.Header().Add(Accept, MediatypeText[MediatypePB])
		w.Header().Add(ContentType, MediatypeText[MediatypePB])
		w.Write(body)
	}
	pm, _ := proto.Marshal(&jsonpb.Simple{OString: proto.String("test")})
	req, _ := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewBuffer(pm))
	req.Header.Add(Accept, MediatypeText[MediatypePB])
	req.Header.Add(ContentType, MediatypeText[MediatypePB])
	resp := httptest.NewRecorder()
	converter(resp, req)
	c := NewPBHTTPConverter()
	m := jsonpb.Simple{}
	// act
	err := c.ReadResponse(resp.Result(), &m)
	// assert
	assert.Nil(t, err)
	assert.Equal(t, "test", m.GetOString())
}

func TestPBHTTPConverter_WriteRequest_Fail(t *testing.T) {
	// arrange
	pm := &jsonpb.Simple{OString: proto.String("test")}
	c := NewPBHTTPConverter()
	// act
	_, err := c.WriteRequest("POST", "http://example.com", -1, pm)
	// assert
	assert.NotNil(t, err)
	assert.Equal(t, "content-type not supported", err.Error())
}

func TestPBHTTPConverter_WriteRequest_JSON(t *testing.T) {
	// arrange
	pm := &jsonpb.Simple{OString: proto.String("test")}
	c := NewPBHTTPConverter()
	// act
	r, err := c.WriteRequest("POST", "http://example.com", MediatypeJSON, pm)
	// assert
	assert.Nil(t, err)
	assert.Equal(t, "POST", r.Method)
	assert.Equal(t, "application/json; charset=UTF-8", r.Header.Get(ContentType))
}

func TestPBHTTPConverter_WriteRequest_PB(t *testing.T) {
	// arrange
	pm := &jsonpb.Simple{OString: proto.String("test")}
	c := NewPBHTTPConverter()
	// act
	r, err := c.WriteRequest("POST", "http://example.com", MediatypePB, pm)
	// assert
	assert.Nil(t, err)
	assert.Equal(t, "POST", r.Method)
	assert.Equal(t, "application/x-protobuf; charset=UTF-8", r.Header.Get(ContentType))
}

func TestPBHTTPConverter_WriteResponse_Fail(t *testing.T) {
	// arrange
	pm := &jsonpb.Simple{OString: proto.String("test")}
	c := NewPBHTTPConverter()
	// act
	_, err := c.WriteResponse(-1, pm)
	// assert
	assert.NotNil(t, err)
	assert.Equal(t, "content-type not supported", err.Error())
}

func TestPBHTTPConverter_WriteBody_JSON(t *testing.T) {
	// arrange
	pm := &jsonpb.Simple{OString: proto.String("test")}
	c := NewPBHTTPConverter()
	// act
	b, err := c.WriteBody(MediatypeJSON, pm)
	// assert
	assert.Nil(t, err)
	assert.Equal(t, `{"o_string":"test"}`, string(b))
}

func TestPBHTTPConverter_WriteBody_PB(t *testing.T) {
	// arrange
	pm := &jsonpb.Simple{OString: proto.String("test")}
	c := NewPBHTTPConverter()
	// act
	b, err := c.WriteBody(MediatypePB, pm)
	// assert
	assert.Nil(t, err)
	assert.Equal(t, "R\x04test", string(b))
}

func TestPBHTTPConverter_WriteResponse_JSON(t *testing.T) {
	// arrange
	pm := &jsonpb.Simple{OString: proto.String("test")}
	c := NewPBHTTPConverter()
	// act
	r, err := c.WriteResponse(MediatypeJSON, pm)
	// assert
	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, "200 OK", r.Status)
	assert.Equal(t, "application/json; charset=UTF-8", r.Header.Get(ContentType))
}

func TestPBHTTPConverter_WriteResponse_PB(t *testing.T) {
	// arrange
	pm := &jsonpb.Simple{OString: proto.String("test")}
	c := NewPBHTTPConverter()
	// act
	r, err := c.WriteResponse(MediatypePB, pm)
	// assert
	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, "200 OK", r.Status)
	assert.Equal(t, "application/x-protobuf; charset=UTF-8", r.Header.Get(ContentType))
}
