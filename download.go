package main

import (
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

type downloadRequest struct {
	Url           string
	Filename      string
	Chunks        int
	ChunksSize    int
	ListChunks    *[][2]int
	TotalSize     int
	SupportRanges bool
	Bar           *pb.ProgressBar
	client        *HttpClient
}

func (d *downloadRequest) Download(idx int, chunk [2]int) error {
	customHeader := map[string]string{
		"Range": fmt.Sprintf("bytes=%v-%v", chunk[0], chunk[1]),
	}
	maps.Copy(customHeader, DefaultHeader)
	resp, err := d.client.NewRequest("GET", d.Url, customHeader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fname := fmt.Sprintf("%v/%v-%v.gsdm.tmp", TmpDir, d.Filename, idx)
	InfoVerbose("Downloading", fname)
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer InfoVerbose(fname, "downloaded")
	proxyReader := d.Bar.NewProxyReader(resp.Body)
	defer f.Close()
	io.Copy(f, proxyReader)
	return nil
}

func (d *downloadRequest) SplitIntoChunks() {
	chunks := make([][2]int, d.Chunks)
	if d.Chunks == 1 {
		chunks[0][0] = 0
		chunks[0][1] = d.TotalSize - 1
		d.ListChunks = &chunks
		return
	}
	for i := range d.Chunks {
		if i == 0 {
			chunks[i][0] = 0
			chunks[i][1] = d.ChunksSize
		} else if i == d.Chunks-1 {
			chunks[i][0] = chunks[i-1][1] + 1
			chunks[i][1] = d.TotalSize - 1
		} else {
			chunks[i][0] = chunks[i-1][1] + 1
			chunks[i][1] = chunks[i][0] + d.ChunksSize
		}
	}
	d.ListChunks = &chunks
}

func (d *downloadRequest) CheckIfTmpFilesExist() bool {
	var sum int64 = 0
	pattern := fmt.Sprintf("%v-*.gsdm.tmp", d.Filename)
	files, err := filepath.Glob(filepath.Join(TmpDir, pattern))
	if err != nil {
		ErrLog(err.Error())
	}
	if len(files) == 0 {
		return false
	}
	size := len(files)
	newChunks := make([][2]int, size)
	for idx := 0; idx < size; idx++ {
		fname := fmt.Sprintf("%v/%v-%v.gsdm.tmp", TmpDir, d.Filename, idx)
		f, err := os.Stat(fname)
		if err != nil || f.IsDir() {
			return false
		}
		newChunks[idx][0] = (*d.ListChunks)[idx][0] + int(f.Size())
		newChunks[idx][1] = (*d.ListChunks)[idx][1]
		sum += f.Size()
	}
	d.Bar.Add64(sum)
	d.ListChunks = &newChunks
	d.Chunks = size
	return true
}

func (d *downloadRequest) MergeDownloadedFile() error {
	InfoVerbose("Start merging tmp files...")
	fname := d.Filename
	if OutputFilaName != "" {
		if err := ValidateDesPath(OutputFilaName); err != nil {
			return err
		}
		fname = OutputFilaName
	}
	out, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer out.Close()
	bar := d.Bar
	bar.SetCurrent(0)
	for idx := 0; idx < d.Chunks; idx++ {
		fname := fmt.Sprintf("%v/%v-%v.gsdm.tmp", TmpDir, d.Filename, idx)
		tmp, err := os.Open(fname)
		if err != nil {
			return err
		}
		defer tmp.Close()
		io.Copy(bar.NewProxyWriter(out), tmp)
	}
	InfoVerbose("file merged successfully:", d.Filename)
	return nil
}

func (d *downloadRequest) Clean() error {
	InfoVerbose("Start cleaning...")
	for idx := 0; idx < d.Chunks; idx++ {
		fname := fmt.Sprintf("%v/%v-%v.gsdm.tmp", TmpDir, d.Filename, idx)
		err := os.Remove(fname)
		if err != nil {
			return fmt.Errorf("can't delete: %v", fname)
		}
	}
	InfoVerbose("Cleaning end")
	return nil
}
