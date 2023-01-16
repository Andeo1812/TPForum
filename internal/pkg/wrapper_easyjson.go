// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package pkg

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

func easyjson3207d587DecodeDbPerformancEprojectInternalPkg(in *jlexer.Lexer, out *ErrResponse) {
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
		case "message":
			out.ErrMassage = string(in.String())
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3207d587EncodeDbPerformancEprojectInternalPkg(out *jwriter.Writer, in ErrResponse) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ErrMassage != "" {
		const prefix string = ",\"message\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ErrMassage))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ErrResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3207d587EncodeDbPerformancEprojectInternalPkg(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ErrResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3207d587EncodeDbPerformancEprojectInternalPkg(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ErrResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3207d587DecodeDbPerformancEprojectInternalPkg(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ErrResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3207d587DecodeDbPerformancEprojectInternalPkg(l, v)
}
