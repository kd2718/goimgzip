package main

import (
	"fmt"
	"github.com/kd2718/goimgzip/imagezip"
	"time"
)

func main() {
	fmt.Println("Heyyyy")

	file := "input/img_urls.txt"

	halt := make(chan int, 20)
	defer close(halt)
	//var utl string
	out, err := imagezip.ReadImageFile(file, halt)

	if err != nil {
		fmt.Println("There was an error", err)
		return
	}

	tex := imagezip.GetImages(out)
	fmt.Println("main done with getimages")

	imagezip.WriteZip(tex)

	<-time.After(30 * time.Second)

	//for {
	//	select{
	//	case utl := <- out:
	//		fmt.Println("test me", utl)
	//	}
	//}

}
