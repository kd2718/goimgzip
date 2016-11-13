package imagezip

import ()

type ByteImage struct {
	file_name string
}

func NewByteImage(reader string) ByteImage {
	bi := ByteImage{
		reader,
	}
	return bi
}
