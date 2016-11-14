package main

import (
	"github.com/kd2718/goimgzip/imagezip"
)

func main() {

	file := "input/img_urls.txt"

	imagezip.StartPipeline(file)


}
