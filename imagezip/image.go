package imagezip

import (
	"fmt"
	"io/ioutil"
	"strings"
	"net/http"
	"sync"
)

func ReadImageFile(filepath string, out_chan, halt chan string) {

	fid, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	lines := strings.Split(string(fid), "\n")

	go func() {
		for _, line := range lines {
			if line == "" {
				return
			}
			out_chan <- line
		}
	}()

}

func GetImages(out chan string) {
	fmt.Println("asdfasdf")

	tex := sync.Mutex{}
	fmt.Println("asdfasdf")
	//var kr string

	for kr := range out {
		fmt.Println("loop", kr)
		go getimage(kr, tex)
	}
}

func getimage(url string, tex sync.Mutex) {
	fmt.Println("test test")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error in get images", err)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("read error", err)
	}

	tex.Lock()
	fmt.Println(resp)
	fmt.Println(resp.Body)
	fmt.Println(data)
	tex.Unlock()
	resp.Body.Close()
}
