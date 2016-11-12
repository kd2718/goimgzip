package imagezip

import (

)
import "net/http"

type ByteImage struct {
	resp *http.Response
}

func NewByteImage(reader *http.Response) ByteImage {
	bi := ByteImage{
		reader,
	}
	return bi
}