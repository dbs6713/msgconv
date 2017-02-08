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
	"encoding/json"
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/valyala/fasthttp"
)

type PBFasthttpConverter struct {
	Charset                int
	MediaTypes             []int
	XPBMessageHeader string
	XPBSchemaHeader  string
}

func NewPBFasthttpConverter() *PBFasthttpConverter {
	return &PBFasthttpConverter{
		Charset:    CharsetUTF8,
		MediaTypes: []int{MediatypeJSON, MediatypePB},
	}
}

func (c *PBFasthttpConverter) CanRead(mediatype int) bool {
	for _, mediaType := range c.MediaTypes {
		if mediaType == mediatype {
			return true
		}
	}
	return false
}

func (c *PBFasthttpConverter) CanWrite(mediatype int) bool {
	return c.CanRead(mediatype)
}

func (c *PBFasthttpConverter) GetMediaType(mediatype string) (int, error) {
	cleanMediaType := CleanMediaType(mediatype)
	for _, mediaType := range c.MediaTypes {
		if MediatypeText[mediaType] == cleanMediaType {
			return mediaType, nil
		}
	}
	return 0, errors.New("media type not supported")
}

func (c *PBFasthttpConverter) ReadRequest(req interface{}, v interface{}) error {
	body := req.(*fasthttp.Request).Body()
	acceptType, err := c.GetMediaType(string(req.(*fasthttp.Request).Header.Peek(Accept)))
	if err != nil {
		return err
	}
	if req.(*fasthttp.Request).Header.IsDelete() || req.(*fasthttp.Request).Header.IsGet() {
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

func (c *PBFasthttpConverter) ReadRequestCtx(ctx *fasthttp.RequestCtx, v interface{}) error {
	body := ctx.Request.Body()
	acceptType, err := c.GetMediaType(string(ctx.Request.Header.Peek(Accept)))
	if err != nil {
		return err
	}
	if ctx.Request.Header.IsDelete() || ctx.Request.Header.IsGet() {
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

func (c *PBFasthttpConverter) ReadResponse(resp interface{}, v interface{}) error {
	body := resp.(*fasthttp.Response).Body()
	contentType, err := c.GetMediaType(string(resp.(*fasthttp.Response).Header.ContentType()))
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

func (c *PBFasthttpConverter) WriteRequest(method string, URI string, contentType int, msg proto.Message) (interface{}, error) {
	var data []byte
	var err error

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	req.Header.Set(CacheControl, "no-cache")
	req.SetRequestURI(URI)

	switch contentType {
	case MediatypeJSON:
		req.Header.Set(Accept, MediatypeText[MediatypeJSON])
		req.Header.Set(ContentType, MediatypeText[MediatypeJSON]+"; charset="+CharsetText[c.Charset])
		data, err = json.Marshal(msg)
		break
	case MediatypePB:
		req.Header.Set(Accept, MediatypeText[MediatypePB])
		req.Header.Set(ContentType, MediatypeText[MediatypePB]+"; charset="+CharsetText[c.Charset])
		data, err = proto.Marshal(msg)
		break
	default:
		return nil, errors.New("unable to marshal message")
	}
	if err != nil {
		return nil, err
	}

	req.Header.Set(ContentLength, string(len(data)))
	req.SetBody(data)
	return req, nil
}

func (c *PBFasthttpConverter) WriteRequestCtx(ctx *fasthttp.RequestCtx, method string, URI string, contentType int, msg proto.Message) error {
	var data []byte
	var err error

	ctx.Request.Header.SetMethod(method)
	ctx.Request.Header.Set(CacheControl, "no-cache")
	ctx.Request.SetRequestURI(URI)

	switch contentType {
	case MediatypeJSON:
		ctx.Request.Header.Set(Accept, MediatypeText[MediatypeJSON])
		ctx.Request.Header.Set(ContentType, MediatypeText[MediatypeJSON]+"; charset="+CharsetText[c.Charset])
		data, err = json.Marshal(msg)
		break
	case MediatypePB:
		ctx.Request.Header.Set(Accept, MediatypeText[MediatypePB])
		ctx.Request.Header.Set(ContentType, MediatypeText[MediatypePB]+"; charset="+CharsetText[c.Charset])
		data, err = proto.Marshal(msg)
		break
	default:
		return errors.New("unable to marshal message")
	}
	if err != nil {
		return err
	}

	ctx.Request.Header.Set(ContentLength, string(len(data)))
	ctx.Request.SetBody(data)
	return nil
}

func (c *PBFasthttpConverter) WriteResponse(contentType int, msg proto.Message) (interface{}, error) {
	var data []byte
	var err error

	resp := fasthttp.AcquireResponse()

	switch contentType {
	case MediatypeJSON:
		resp.Header.Set(Accept, MediatypeText[MediatypeJSON])
		resp.Header.Set(ContentType, MediatypeText[MediatypeJSON]+"; charset="+CharsetText[c.Charset])
		data, err = json.Marshal(msg)
		break
	case MediatypePB:
		resp.Header.Set(Accept, MediatypeText[MediatypePB])
		resp.Header.Set(ContentType, MediatypeText[MediatypePB]+"; charset="+CharsetText[c.Charset])
		data, err = proto.Marshal(msg)
		break
	default:
		return nil, errors.New("unable to marshal message")
	}
	if err != nil {
		return nil, err
	}

	resp.Header.Set(ContentLength, string(len(data)))
	resp.SetBody(data)
	return resp, nil
}

func (c *PBFasthttpConverter) WriteResponseCtx(ctx *fasthttp.RequestCtx, contentType int, msg proto.Message) error {
	var data []byte
	var err error

	switch contentType {
	case MediatypeJSON:
		ctx.Response.Header.Set(Accept, MediatypeText[MediatypeJSON])
		ctx.Response.Header.Set(ContentType, MediatypeText[MediatypeJSON]+"; charset="+CharsetText[c.Charset])
		data, err = json.Marshal(msg)
		break
	case MediatypePB:
		ctx.Response.Header.Set(Accept, MediatypeText[MediatypePB])
		ctx.Response.Header.Set(ContentType, MediatypeText[MediatypePB]+"; charset="+CharsetText[c.Charset])
		data, err = proto.Marshal(msg)
		break
	default:
		return errors.New("unable to marshal message")
	}
	if err != nil {
		return err
	}

	ctx.Response.Header.Set(ContentLength, string(len(data)))
	ctx.Response.SetBody(data)
	return nil
}
