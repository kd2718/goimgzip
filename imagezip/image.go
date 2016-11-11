package imagezip

import "io"

type ByteImage struct {
	io.Reader
}

func NewByteImage(reader io.Reader) ByteImage {
	bi := ByteImage{
		reader,
	}
	return bi
}