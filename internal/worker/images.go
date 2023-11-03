package worker

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

func getImages(url, hdUrl string) ([]byte, []byte, error) {
	const op = "worker.getImages"

	var wg sync.WaitGroup
	imageCh := make(chan []byte)
	hdImageCh := make(chan []byte)
	var imageErr, hdImageErr error

	wg.Add(1)
	go getImage(url, &wg, imageCh, &imageErr)

	wg.Add(1)
	go getImage(hdUrl, &wg, hdImageCh, &hdImageErr)

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

func getImage(url string, wg *sync.WaitGroup, ch chan []byte, returnedErr *error) {
	defer wg.Done()

	imageData, err := downloadImage(url)
	if err != nil {
		*returnedErr = err
		return
	}

	ch <- imageData
}

func downloadImage(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}
