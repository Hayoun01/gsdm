package main

import (
	"errors"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	tempFileSuffix = ".gsdm.tmp"
	tempFileTTL    = 400 * time.Second
)

func isNotValidFilename(filename string) bool {
	filename = filepath.Base(path.Clean(filename))
	return filename == "" || filename == "." || filename == "/"
}

func ValidateDesPath(path string) error {
	fullPath := filepath.Dir(path)
	dir, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
				return fmt.Errorf("error making new directory: %v", err)
			} else {
				return nil
			}
		} else {
			return fmt.Errorf("error checking Dest directory: %v", err)
		}
	}
	if !dir.IsDir() {
		return errors.New(fullPath + ": Not a directory")
	}
	return nil
}

func GuessFileName(resp *http.Response) (string, error) {
	filename := filepath.Base(path.Clean(OutputFilaName))
	if !isNotValidFilename(filename) {
		return filename, nil
	}
	filename = resp.Request.URL.Path
	if disp := resp.Header.Get("Content-Disposition"); disp != "" {
		if _, parsedFileName, err := mime.ParseMediaType(disp); err == nil {
			if val, ok := parsedFileName["filename"]; ok {
				filename = val
			}
		}

	}
	filename = filepath.Base(path.Clean(filename))
	if isNotValidFilename(filename) {
		return "", errors.New("filename couldn't be determined")
	}
	return filename, nil
}

func CleanTmpDir() error {
	InfoVerbose("Start checking and cleaning temporary directory....")
	now := time.Now()
	tempDir := os.TempDir()
	pattern := "*" + tempFileSuffix
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return errors.New("error reading temporary directory: " + err.Error())
	}
	for _, file := range files {
		if !file.IsDir() {
			fileInfo, err := file.Info()
			if err != nil {
				continue
			}
			ok, _ := filepath.Match(pattern, fileInfo.Name())
			if ok {
				if now.Sub(fileInfo.ModTime()) > tempFileTTL {
					err := os.Remove(filepath.Join(tempDir, file.Name()))
					if err != nil {
						return errors.New("error deleting temporary file: " + err.Error())
					}
					InfoVerbose(file.Name(), "deleted")
				}
			}
		}
	}
	return nil
}

func InfoLog(args ...string) {
	log.Println(LogsColorize(ColorBlue, "INFO:", args...))
}

func WarnLog(args ...string) {
	log.Println(LogsColorize(ColorYellow, "WARNING:", args...))
}

func ErrLog(args ...string) {
	log.Println(LogsColorize(ColorRed, "ERROR:", args...))
	os.Exit(-1)
}

func InfoVerbose(args ...string) {
	if Verbose {
		fmt.Printf("\n")
		log.Print(LogsColorize(ColorBlue, "INFO:", args...))
	}
}
