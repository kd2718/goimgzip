package imagezip

import (
	"compress/gzip"
	"errors"
	"fmt"
	//"io/ioutil"
	"net/http"
	"os"
	"bufio"
	"io/ioutil"
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
	 	line, _,  err := lines.ReadLine()
		 strLine := string(line)
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
				fmt.Println("myline", strLine)
				if strLine == "" {
					return
				}
				out_chan <- strLine
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
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println("error in get images", err)
					return
				}
				fmt.Println("stats", resp.Status)
				//defer resp.Body.Close()
				blobImage := NewByteImage(resp)

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

	my_file, err := os.Create("out/test.zip")
	if err != nil {
		fmt.Println("good stuff... not", err)
	}
	defer my_file.Close()
	w, _ := gzip.NewWriterLevel(my_file, gzip.NoCompression)
	defer w.Close()
	idx := 1

	go func() {
		for blob := range tex {
			go func() {
				fid, err := ioutil.TempFile("out/", "image_")

				defer os.Remove(fid.Name())
				defer fid.Close()
				car := io.Rea(blob.resp.Body)

				copied, err := w.Write(car)

				//copied, err := io.Copy(w, blob.resp.Body)
				if err != nil {
					fmt.Println("err is not nill!!!", err)
				}
				fmt.Println("copied:", copied)
				//w.
				defer blob.resp.Body.Close()

				w.Flush()

				idx++
			}()
		}
	}()
}

