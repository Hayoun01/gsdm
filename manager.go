package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/cheggaaa/pb/v3"
	"gitlab.com/poldi1405/go-ansi"
)

func Run(url string) error {
	client := NewHttpClient()
	resp, err := client.NewRequest("HEAD", url, DefaultHeader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fname, err := GuessFileName(resp)
	if err != nil {
		return err
	}
	sizeFileByBytes, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}
	d := downloadRequest{
		Url:           url,
		Filename:      fname,
		TotalSize:     sizeFileByBytes,
		SupportRanges: resp.Header.Get("Accept-Ranges") == "bytes",
		client:        client,
	}
	d.Bar = NewBar(d.TotalSize)
	if d.SupportRanges {
		d.Chunks = Workers
		d.ChunksSize = sizeFileByBytes / d.Chunks
		d.SplitIntoChunks()
		if ok := d.CheckIfTmpFilesExist(); ok {
			InfoLog("File found completing download...")
		}
	} else {
		WarnLog("Server does not support byte range requests. Downloading the whole file sequentially...")
		d.Chunks = 1
		d.ChunksSize = d.TotalSize
		d.SplitIntoChunks()
	}
	errChan := make(chan error, len(*d.ListChunks))
	var wg sync.WaitGroup
	for idx, chunk := range *d.ListChunks {
		wg.Add(1)
		go func(idx int, chunk [2]int) {
			defer wg.Done()
			if err := d.Download(idx, chunk); err != nil {
				errChan <- err
			}
		}(idx, chunk)
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		return err
	}
	if err := d.MergeDownloadedFile(); err != nil {
		return err
	}
	d.Bar.Finish()
	fmt.Print("\033[A", ansi.ClearLine())
	fmt.Println("✔", Colorize(ColorGreen, d.Filename))
	if err := d.Clean(); err != nil {
		return err
	}
	return nil
}

func NewBar(totalSize int) *pb.ProgressBar {
	bar := pb.Full.Start(totalSize)
	bar.SetWidth(120)
	bar.SetTemplateString(`{{percent . }}{{bar . "▕" "\033[34m█\033[0m" "\033[34m█\033[0m" " " "▏" }}{{counters . "%7s/%7s"}} @ {{speed . "%s/s" | green }}`)
	bar.Set(pb.SIBytesPrefix, true)
	return bar
}
