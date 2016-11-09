package main

import (
	"fmt"
	"github.com/kd2718/goimgzip/imagezip"
	"time"
)

func main() {
	fmt.Println("Heyyyy")

	file := "/Users/koryd/img_urls.txt"

	out := make(chan string, 20)
	halt := make(chan string, 20)
	//var utl string
	imagezip.ReadImageFile(file, out, halt)

	go imagezip.GetImages(out)

	<-time.After(10 * time.Second)

	//for {
	//	select{
	//	case utl := <- out:
	//		fmt.Println("test me", utl)
	//	}
	//}

}
