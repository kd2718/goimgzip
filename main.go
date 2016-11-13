package main

import (
	"fmt"
	"github.com/kd2718/goimgzip/imagezip"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	fmt.Println("Heyyyy")

	file := "input/img_urls.txt"
	var wg sync.WaitGroup
	halt := make(chan int, 20)
	defer close(halt)
	//var utl string
	wg.Add(1)
	out, err := imagezip.ReadImageFile(file, halt, &wg)

	if err != nil {
		fmt.Println("There was an error", err)
		return
	}

	tex := imagezip.GetImages(out)
	_ = tex
	fmt.Println("main done with getimages")

	tw := imagezip.WriteZip(tex, &wg)
	defer tw.Close()
	wg.Wait()
	end := time.Since(start)
	total := time.Duration(end)
	fmt.Println("total time:", total, end)

}
