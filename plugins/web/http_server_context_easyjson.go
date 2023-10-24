/*
 *  Copyright (c) 2022-2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package web

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonF7d450f8DecodeGoOsspkgComGoppyPluginsWeb(in *jlexer.Lexer, out *errMessage) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "msg":
			out.Message = string(in.String())
		case "ctx":
			(out.Ctx).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF7d450f8EncodeGoOsspkgComGoppyPluginsWeb(out *jwriter.Writer, in errMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"msg\":"
		out.RawString(prefix[1:])
		out.String(string(in.Message))
	}
	if len(in.Ctx) != 0 {
		const prefix string = ",\"ctx\":"
		out.RawString(prefix)
		(in.Ctx).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v errMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF7d450f8EncodeGoOsspkgComGoppyPluginsWeb(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v errMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF7d450f8EncodeGoOsspkgComGoppyPluginsWeb(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *errMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF7d450f8DecodeGoOsspkgComGoppyPluginsWeb(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *errMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF7d450f8DecodeGoOsspkgComGoppyPluginsWeb(l, v)
}
func easyjsonF7d450f8DecodeGoOsspkgComGoppyPluginsWeb1(in *jlexer.Lexer, out *ErrCtx) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
	} else {
		in.Delim('{')
		*out = make(ErrCtx)
		for !in.IsDelim('}') {
			key := string(in.String())
			in.WantColon()
			var v1 interface{}
			if m, ok := v1.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := v1.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				v1 = in.Interface()
			}
			(*out)[key] = v1
			in.WantComma()
		}
		in.Delim('}')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonF7d450f8EncodeGoOsspkgComGoppyPluginsWeb1(out *jwriter.Writer, in ErrCtx) {
	if in == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
		out.RawString(`null`)
	} else {
		out.RawByte('{')
		v2First := true
		for v2Name, v2Value := range in {
			if v2First {
				v2First = false
			} else {
				out.RawByte(',')
			}
			out.String(string(v2Name))
			out.RawByte(':')
			if m, ok := v2Value.(easyjson.Marshaler); ok {
				m.MarshalEasyJSON(out)
			} else if m, ok := v2Value.(json.Marshaler); ok {
				out.Raw(m.MarshalJSON())
			} else {
				out.Raw(json.Marshal(v2Value))
			}
		}
		out.RawByte('}')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v ErrCtx) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF7d450f8EncodeGoOsspkgComGoppyPluginsWeb1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ErrCtx) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF7d450f8EncodeGoOsspkgComGoppyPluginsWeb1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ErrCtx) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF7d450f8DecodeGoOsspkgComGoppyPluginsWeb1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ErrCtx) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF7d450f8DecodeGoOsspkgComGoppyPluginsWeb1(l, v)
}
