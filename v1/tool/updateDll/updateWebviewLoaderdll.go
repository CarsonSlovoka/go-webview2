package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	const zipPath = "webview.zip"
	const downloadURL = "https://globalcdn.nuget.org/packages/microsoft.web.webview2.1.0.1661.34.nupkg"
	if err := download(downloadURL, zipPath); err != nil {
		log.Fatal(err)
	}

	const unZipPath = "./tempWebview"
	if err := UnZip(zipPath, unZipPath); err != nil {
		log.Fatal(err)
	}
}

func download(src, dst string) error {
	log.Printf("download from %q to %q", src, dst)
	resp, err := http.Get(src)
	if err != nil {
		return err
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if _, err = f.Write(contents); err != nil {
		return err
	}
	return nil
}

func UnZip(zipPath, dstDirPath string) error {
	// open zip
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = r.Close()
	}()

	var (
		dst *os.File
		src io.ReadCloser
	)

	for _, zipFile := range r.File {
		if zipFile.FileInfo().IsDir() {
			if err = os.MkdirAll(filepath.Join(dstDirPath, zipFile.Name), os.ModePerm); err != nil {
				log.Fatal(err)
			}
			continue
		}

		// prepare: dst
		dstFilePath := filepath.Join(dstDirPath, zipFile.Name)
		dst, err = os.Create(dstFilePath)
		if err != nil {
			log.Fatal(err)
		}

		// prepare: src
		src, err = zipFile.Open()
		if err != nil {
			return err
		}

		// copy from src to dst
		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}
		_ = dst.Close()
		_ = src.Close()

		if filepath.Base(zipFile.Name) == "WebView2Loader.dll" {
			// absDstFilePath, _ := filepath.Abs(dstFilePath)
			fmt.Println(zipFile.Name)
		}
	}
	return nil
}
