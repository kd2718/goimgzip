package imagezip

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	BYTE_IMAGE_CHAN_SIZE = 100
	URL_CHAN_SIZE = 200
)

func ReadImageFile(filepath string, wg *sync.WaitGroup) (chan string, error) {

	out_chan := make(chan string, URL_CHAN_SIZE)

	go readImage(filepath, out_chan, wg)

	return out_chan, nil

}

func readImage(filepath string, out_chan chan string, wg *sync.WaitGroup) {
	fid, err := os.Open(filepath)
	defer fid.Close()
	if err != nil {
		fmt.Println("error", err)
		return
	}
	lines := bufio.NewReader(fid)
	for {
		line, _, err := lines.ReadLine()
		strLine := string(line)
		if err != nil {
			fmt.Println("end of file...")
			break
		}
		if strLine == "" {
			return
		}
		wg.Add(1)
		out_chan <- strLine
	}
}

func GetImages(out chan string) chan ByteImage {
	tex := make(chan ByteImage, BYTE_IMAGE_CHAN_SIZE)
	go getImagesLoop(out, tex)
	return tex
}

func getImagesLoop(out chan string, tex chan ByteImage) {
	for url := range out {
		go getImage(url, tex)
	}
}

func getImage(url string, tex chan ByteImage) {
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
}

func WriteZip(tex chan ByteImage, wg *sync.WaitGroup) *zip.Writer {
	fmt.Println("time to write the zip")
	wg.Done()

	my_file, err := os.Create("out/test.zip")
	if err != nil {
		fmt.Println("good stuff... not", err)
	}
	tw := zip.NewWriter(my_file)

	go writeBlobtoZip(tw, tex, wg)
	return tw
}

func writeBlobtoZip(tw *zip.Writer, tex chan ByteImage, wg *sync.WaitGroup) {
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

	}
}


func StartPipeline(file string){
	start := time.Now()
	var wg sync.WaitGroup
	//var utl string
	wg.Add(1)
	out, err := ReadImageFile(file, &wg)

	if err != nil {
		fmt.Println("There was an error", err)
		return
	}

	tex := GetImages(out)
	fmt.Println("main done with getimages")

	tw := WriteZip(tex, &wg)
	defer tw.Close()
	wg.Wait()
	end := time.Since(start)
	total := time.Duration(end)
	fmt.Println("total time:", total, end)
}