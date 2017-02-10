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
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
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

func (p *PBHTTPConverter)ReadRequest(req *http.Request, msg interface{}) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	acceptType, err := p.GetMediaType(string(req.Header.Get(Accept)))
	if err != nil {
		return err
	}
	if req.Method == http.MethodDelete || req.Method == http.MethodGet {
		return nil
	}

	switch acceptType {
	case MediatypeJSON:
		err = json.Unmarshal(body, &msg)
		break
	case MediatypePB:
		err = proto.Unmarshal(body, msg.(proto.Message))
		break
	default:
		return errors.New("unable to unmarshal message")
	}
	if err != nil {
		return err
	}
	return nil
}

func (p *PBHTTPConverter)ReadResponse(resp *http.Response, msg interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	contentType, err := p.GetMediaType(resp.Header.Get(ContentType))
	if err != nil {
		return err
	}

	switch contentType {
	case MediatypeJSON:
		err = json.Unmarshal(body, &msg)
		break
	case MediatypePB:
		err = proto.Unmarshal(body, msg.(proto.Message))
		break
	default:
		return errors.New("content-type not supported")
	}
	if err != nil {
		return err
	}
	return nil
}

func (p *PBHTTPConverter)WriteRequest(method string, URL string, contentType int, msg proto.Message) (*http.Request, error) {
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
		return nil, errors.New("content-type not supported")
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

func (p *PBHTTPConverter)WriteResponse(contentType int, msg proto.Message) (*http.Response, error) {
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
		return nil, errors.New("content-type not supported")
	}
	if err != nil {
		return nil, err
	}
	resp := http.Response{
		Status: "200 OK",
		StatusCode: 200,
		Proto: "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	resp.Header = http.Header{}
	resp.Header.Set(Accept, mediaType)
	resp.Header.Set(ContentType, mediaType+"; charset="+CharsetText[p.Charset])
	resp.Header.Set(ContentLength, string(len(data)))
	r := bytes.NewReader(data)
	rc, ok := io.Reader(r).(io.ReadCloser)
	if !ok && r != nil {
		rc = ioutil.NopCloser(rc)
	}
	resp.Body = rc
	return &resp, nil
}
