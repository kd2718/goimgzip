package imagezip

import (
	"compress/gzip"
	"errors"
	"fmt"
	//"io/ioutil"
	"net/http"
	"os"
	"bufio"
	"strconv"
	"io"
)

func ReadImageFile(filepath string, halt chan int) (chan string, error) {

	out_chan := make(chan string, 20)

	fid, err := os.Open(filepath)
	if err != nil {
		fmt.Println("error", err)
		return out_chan, err
	}
	lines := bufio.NewReader(fid)

	//lines := strings.Split(string(fid), "\n")

	for {
	 	line, err := lines.ReadString('\n')
		if err != nil {
			fmt.Println("end of file...")
			break
		}
		select {
		case <-halt:
			fmt.Println("Halt detected")
			return out_chan, errors.New("Halted...")
		default:
			fmt.Println("Get the good stuff")
			go func() {
				fmt.Println("myline", line)
				if line == "" {
					return
				}
				out_chan <- line
				return
			}()
		}
	}

	return out_chan, err

}

func GetImages(out chan string) chan ByteImage {
	tex := make(chan ByteImage, 10)
	go func() {
		for url := range out {
			go func() {
				fmt.Println("test test")
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println("error in get images", err)
					return
				}
				//data, err := ioutil.ReadAll(resp.Body)

				blobImage := NewByteImage(resp.Body)

				if err != nil {
					fmt.Println("read error", err)
				}
				fmt.Println("Writing to chan")
				tex <- blobImage
				return
			}()
		}
	}()
	return tex
}

func WriteZip(tex chan ByteImage) {
	fmt.Println("time to write the zip")

	my_file, err := os.Create("test.zip")
	if err != nil {
		fmt.Println("good stuff... not", err)
	}
	defer my_file.Close()
	w, _ := gzip.NewWriterLevel(my_file, gzip.NoCompression)
	defer w.Close()
	idx := 1
	var af []byte

	go func(){
		for blob := range tex {
			fmt.Println("starting blob write")
			raw_file, _ := os.Create(strconv.Itoa(idx))
			_ = raw_file
			blob.Read(af)
			io.Copy(w, blob.Reader)
			//if err != nil {
			//	fmt.Println("Error now...", err)
			//}
			//
			//_, err = fid.Write(blob.bytes)
			//
			//if err != nil {
			//	fmt.Println("err over here now", err)
			//}
			w.Flush()

			idx++
		}
	}()
}

