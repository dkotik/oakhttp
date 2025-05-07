package oakhttp

import (
	"encoding/json"
	"fmt"
	"io"
)

type Codec interface {
	Encode(io.Writer, any) error
	Decode(any, io.Reader) error
}

var codecs = map[string]Codec{
	"application/json": &JSONCodec{},
}

func GetCodec(contentType string) (Codec, error) {
	c, ok := codecs[contentType]
	if !ok {
		return nil, NewInvalidRequestError(fmt.Errorf(
			"there is no known codec for content type %q", contentType))
	}
	return c, nil
}

func SetCodec(contentType string, c Codec) {
	if c == nil {
		panic("cannot set <nil> codec")
	}
	codecs[contentType] = c
}

type JSONCodec struct{}

func (j *JSONCodec) Encode(w io.Writer, data any) error {
	return json.NewEncoder(w).Encode(data)
}

func (j *JSONCodec) Decode(data any, r io.Reader) error {
	return json.NewDecoder(r).Decode(&data)
}
