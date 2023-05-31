package ws

type ICodec interface {

	// []byte -> struct
	Encode(b []byte) (int, any)

	// struct -> []byte
	Decode(t int, v any) []byte
}
