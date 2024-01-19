// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package model

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

func easyjsonB4c90860DecodeAssetsCommonModel(in *jlexer.Lexer, out *Sort) {
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
		case "field":
			out.Field = string(in.String())
		case "order":
			out.Order = SortOrder(in.String())
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
func easyjsonB4c90860EncodeAssetsCommonModel(out *jwriter.Writer, in Sort) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Field != "" {
		const prefix string = ",\"field\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Field))
	}
	if in.Order != "" {
		const prefix string = ",\"order\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Order))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Sort) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB4c90860EncodeAssetsCommonModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Sort) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB4c90860EncodeAssetsCommonModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Sort) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB4c90860DecodeAssetsCommonModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Sort) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB4c90860DecodeAssetsCommonModel(l, v)
}