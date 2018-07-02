package jwalk

import (
	"encoding/json"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"github.com/pkg/errors"
)

// ObjectWalker iterates through JSON object fields.
type ObjectWalker interface {
	Walk(func(name string, value interface{}) error) error
	json.Marshaler
}

// object represents JSON object.
type object struct {
	fields []field
}

// field represents JSON object field.
type field struct {
	name  string
	value interface{}
}

func (o object) Walk(fn func(name string, value interface{}) error) error {
	for _, f := range o.fields {
		if err := fn(f.name, f.value); err != nil {
			return err
		}
	}

	return nil
}

func (o object) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	if err := o.marshal(&w); err != nil {
		return nil, errors.Wrap(err, "marshal object")
	}
	return w.Buffer.BuildBytes(), w.Error
}

func (o object) marshal(w *jwriter.Writer) error {
	w.RawByte('{')
	first := true
	for _, f := range o.fields {
		prefix := ",\"" + f.name + "\":"
		if first {
			first = false
			w.RawString(prefix[1:])
		} else {
			w.RawString(prefix)
		}

		mlr, ok := f.value.(json.Marshaler)
		if !ok {
			return errors.New("failed to assert value to json.Marshaler")
		}
		data, err := mlr.MarshalJSON()
		if err != nil {
			return errors.Wrap(err, "marshal object iterator")
		}
		w.Raw(data, nil)
	}
	w.RawByte('}')
	return w.Error
}

// ObjectsWalker iterates through array of JSON objects.
type ObjectsWalker interface {
	Walk(func(obj ObjectWalker) error) error
	json.Marshaler
}

type objects []ObjectWalker

func (o objects) Walk(fn func(obj ObjectWalker) error) error {
	for _, obj := range o {
		if err := fn(obj); err != nil {
			return err
		}
	}

	return nil
}

func (o objects) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	if err := o.marshal(&w); err != nil {
		return nil, errors.Wrap(err, "marshal objects")
	}
	return w.Buffer.BuildBytes(), w.Error
}

func (o objects) marshal(w *jwriter.Writer) error {
	w.RawByte('[')
	first := true
	for _, obj := range o {
		if first {
			first = false
		} else {
			w.RawByte(',')
		}
		data, err := obj.MarshalJSON()
		if err != nil {
			return errors.Wrap(err, "marshal object iterator")
		}
		w.Raw(data, nil)
	}
	w.RawByte(']')
	return w.Error
}

// Parse parses the json data.
func Parse(data []byte) (interface{}, error) {
	l := jlexer.Lexer{Data: data}
	return parse(&l)
}

func parse(l *jlexer.Lexer) (interface{}, error) {
	switch {
	case l.IsDelim('{'):
		return getObject(l)
	case l.IsDelim('['):
		return getArray(l)
	default:
		return getValue(l)
	}
}

func getObject(l *jlexer.Lexer) (ObjectWalker, error) {
	var obj object
	l.Delim('{')
	for !l.IsDelim('}') {
		key := l.UnsafeString()
		field := field{name: key}
		l.WantColon()
		if l.IsNull() {
			field.value = value{[]byte("null")}
			obj.fields = append(obj.fields, field)
			l.Skip()
			l.WantComma()
			continue
		}

		val, err := parse(l)
		if err != nil {
			return obj, errors.Wrap(err, "get object")
		}

		field.value = val
		obj.fields = append(obj.fields, field)
		l.WantComma()
	}
	l.Delim('}')

	return obj, l.Error()
}

func getArray(l *jlexer.Lexer) (interface{}, error) {
	raw := l.Raw()
	rc := make([]byte, len(raw))
	copy(rc, raw)
	ll := jlexer.Lexer{Data: rc}
	ll.Delim('[')
	defer ll.Delim(']')
	switch {
	case ll.IsDelim('{'):
		objs := make(objects, 0)
		for !ll.IsDelim(']') {
			obj, err := getObject(&ll)
			if err != nil {
				return objs, errors.Wrap(err, "array of objects")
			}
			objs = append(objs, obj)
			ll.WantComma()
		}
		return objs, l.Error()
	default:
		rc := make([]byte, len(raw))
		copy(rc, raw)
		ll = jlexer.Lexer{Data: rc}
		return getValue(&ll)
	}
}

func getValue(l *jlexer.Lexer) (Value, error) {
	raw := l.Raw()
	return value{raw}, l.Error()
}
