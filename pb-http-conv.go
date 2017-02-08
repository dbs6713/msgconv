// Copyright 2017 Intellectual Reserve, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.
//
// @author Don B. Stringham <stringhamdb@ldschurch.org>
//
package msgconv

import (
	"bytes"
	"errors"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
	"io"
)

type PBHTTPConverter struct {
	Charset                int
	MediaTypes             []int
	XPBMessageHeader string
	XPBSchemaHeader  string
}

func NewPBHTTPConverter() *PBHTTPConverter {
	return &PBHTTPConverter{
		Charset:    CharsetUTF8,
		MediaTypes: []int{MediatypeJSON, MediatypePB},
	}
}

func (p *PBHTTPConverter)CanRead(mediatype int) bool {
	for _, mediaType := range p.MediaTypes {
		if mediaType == mediatype {
			return true
		}
	}
	return false
}

func (p *PBHTTPConverter)CanWrite(mediatype int) bool {
	return p.CanRead(mediatype)
}

func (p *PBHTTPConverter)GetMediaType(mediatype string) (int, error) {
	cleanMediaType := CleanMediaType(mediatype)
	for _, mediaType := range p.MediaTypes {
		if MediatypeText[mediaType] == cleanMediaType {
			return mediaType, nil
		}
	}
	return 0, errors.New("media type not supported")
}

func (p *PBHTTPConverter)ReadRequest(req interface{}, v interface{}) error {
	body, err := ioutil.ReadAll(req.(http.Request).Body)
	if err != nil {
		return err
	}
	acceptType, err := p.GetMediaType(string(req.(http.Request).Header.Get(Accept)))
	if err != nil {
		return err
	}
	if req.(http.Request).Method == http.MethodDelete || req.(http.Request).Method == http.MethodGet {
		return nil
	}

	switch acceptType {
	case MediatypeJSON:
		err = json.Unmarshal(body, &v)
		break
	case MediatypePB:
		err = proto.Unmarshal(body, v.(proto.Message))
		break
	default:
		return errors.New("unable to unmarshal message")
	}
	if err != nil {
		return err
	}
	return nil
}

func (p *PBHTTPConverter)ReadResponse(resp interface{}, v interface{}) error {
	body, err := ioutil.ReadAll(resp.(http.Response).Body)
	contentType, err := p.GetMediaType(resp.(http.Response).Header.Get(ContentType))
	if err != nil {
		return err
	}

	switch contentType {
	case MediatypeJSON:
		err = json.Unmarshal(body, &v)
		break
	case MediatypePB:
		err = proto.Unmarshal(body, v.(proto.Message))
		break
	default:
		return errors.New("unable to unmarshal message")
	}
	if err != nil {
		return err
	}
	return nil
}

func (p *PBHTTPConverter)WriteRequest(method string, URL string, contentType int, msg proto.Message) (interface{}, error) {
	var data []byte
	var err error
	var mediaType string

	switch contentType {
	case MediatypeJSON:
		mediaType = MediatypeText[MediatypeJSON]
		data, err = json.Marshal(msg)
		break
	case MediatypePB:
		mediaType = MediatypeText[MediatypePB]
		data, err = proto.Marshal(msg)
		break
	default:
		return nil, errors.New("unable to marshal message")
	}
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, URL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set(Accept, mediaType)
	req.Header.Set(ContentType, mediaType+"; charset="+CharsetText[p.Charset])
	req.Header.Set(CacheControl, "no-cache")
	req.Header.Set(ContentLength, string(len(data)))
	return req, nil
}

func (p *PBHTTPConverter)WriteResponse(contentType int, msg proto.Message) (interface{}, error) {
	var data []byte
	var err error
	var mediaType string

	switch contentType {
	case MediatypeJSON:
		mediaType = MediatypeText[MediatypeJSON]
		data, err = json.Marshal(msg)
		break
	case MediatypePB:
		mediaType = MediatypeText[MediatypePB]
		data, err = proto.Marshal(msg)
		break
	default:
		return nil, errors.New("unable to marshal message")
	}
	if err != nil {
		return nil, err
	}
	resp := http.Response{}
	resp.Header.Set(Accept, mediaType)
	resp.Header.Set(ContentType, mediaType+"; charset="+CharsetText[p.Charset])
	resp.Header.Set(ContentLength, string(len(data)))
	r := bytes.NewReader(data)
	rc, ok := io.Reader(r).(io.ReadCloser)
	if !ok && r != nil {
		rc = ioutil.NopCloser(rc)
	}
	resp.Body = rc
	return resp, nil
}
