package worker

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func getImages(url, hdUrl string) (*bytes.Buffer, *bytes.Buffer, error) {
	const op = "worker.getImages"

	var wg sync.WaitGroup
	imageCh := make(chan *bytes.Buffer)
	hdImageCh := make(chan *bytes.Buffer)
	var imageErr, hdImageErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := getImage(url, imageCh)
		if err != nil {
			imageErr = err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := getImage(hdUrl, hdImageCh)
		if err != nil {
			hdImageErr = err
		}
	}()

	go func() {
		wg.Wait()
		close(imageCh)
		close(hdImageCh)
	}()

	imageData := <-imageCh
	hdImageData := <-hdImageCh

	if imageErr != nil || hdImageErr != nil {
		return nil, nil, fmt.Errorf("%s: %v, %v", op, imageErr, hdImageErr)
	}

	return imageData, hdImageData, nil
}

func getImage(url string, ch chan *bytes.Buffer) error {
	imageData, err := downloadImage(url)
	if err != nil {
		return err
	}

	ch <- imageData

	return nil
}

func downloadImage(url string) (*bytes.Buffer, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	imageData := &bytes.Buffer{}

	_, err = io.Copy(imageData, response.Body)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}
