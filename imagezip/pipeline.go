package imagezip

import (
	_ "compress/gzip"
	"errors"
	"fmt"
	//"io/ioutil"
	_ "archive/tar"
	"bufio"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"archive/zip"
	"sync"
)

func ReadImageFile(filepath string, halt chan int, wg *sync.WaitGroup) (chan string, error) {

	out_chan := make(chan string, 20)

	fid, err := os.Open(filepath)
	if err != nil {
		fmt.Println("error", err)
		return out_chan, err
	}
	lines := bufio.NewReader(fid)

	//lines := strings.Split(string(fid), "\n")

	for {
		line, _, err := lines.ReadLine()
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
			go func() {
				if strLine == "" {
					return
				}
				wg.Add(1)
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
			go func(url string) {
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println("error in get images", err)
					return
				}
				fid, err := ioutil.TempFile("out/", "image_")
				defer fid.Close()
				if err != nil {
					fmt.Println("err creating temp file", err)
				}

				fmt.Println("stats", resp.Status)
				//defer resp.Body.Close()
				_, err = io.Copy(fid, resp.Body)
				if err != nil {
					fmt.Println("error copy body to temp file", err)
					os.Remove(fid.Name())
				}
				blobImage := NewByteImage(fid.Name())

				if err != nil {
					fmt.Println("read error", err)
				}
				tex <- blobImage
				return
			}(url)
		}
	}()
	return tex
}

func WriteZip(tex chan ByteImage, wg *sync.WaitGroup) *zip.Writer {
	fmt.Println("time to write the zip")
	wg.Done()

	my_file, err := os.Create("out/test.zip")
	if err != nil {
		fmt.Println("good stuff... not", err)
	}
	//tw := tar.NewWriter(my_file)
	tw := zip.NewWriter(my_file)

	//defer my_file.Close()

	go func() {
		for blob := range tex {

			//defer blob.resp.Close()

			myfile, err := os.Open(blob.file_name)
			if err != nil {
				fmt.Println("errer open tempfile", err)
			}

			wr, err := tw.Create(blob.file_name)
			if err != nil {
				fmt.Println("error creating thing", err)
			}

			io.Copy(wr, myfile)
			wg.Done()
			myfile.Close()
			os.Remove(myfile.Name())



			//fstat, err := myfile.Stat()
			//if err != nil {
			//	fmt.Println("there was an error", err)
			//}
			//header, err := tar.FileInfoHeader(fstat, blob.file_name)
			//if err != nil {
			//	fmt.Println("err or in fileinfoheardr", err)
			//}
			//
			//fmt.Println("File header info:", header)
			//
			//tw.WriteHeader(header)
			//temp, err := ioutil.ReadAll(blob.resp)
			//if err != nil {
			//	fmt.Println("error reading file", err)
			//}
			//
			//_, err = tw.Write(temp)
			//if err != nil {
			//	fmt.Println("error write", err)
			//}

			//io.Copy(tw, myfile)
			//tw.Flush()

		}
	}()
	return tw
}
