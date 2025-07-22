package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CreateZipArchive(config map[string]interface{}, archiveFile string) error {
	zipFile, err := os.Create(archiveFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWritter := zip.NewWriter(zipFile)
	defer zipWritter.Close()

	targets := config["targets"].([]interface{})

	for _, target := range targets {
		var path string
		var exclude string

		if targetStr, ok := target.(string); ok {
			path = targetStr
		} else {
			targetMap := target.(map[string]interface{})
			path = targetMap["path"].(string)
			if excludeVal, exists := targetMap["exclude"]; exists {
				exclude = excludeVal.(string)
			}
		}

		err := AddObjectsToZip(zipWritter, path, exclude)
		if err != nil {
			return err
		}
	}

	fmt.Printf("ZIP архив создан: %s\n", archiveFile)
	return nil
}

func AddObjectsToZip(zipWritter *zip.Writer, path, exclude string) error {
	files, err := filepath.Glob(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if exclude != "" && matchesExclude(file, exclude) {
			continue
		}

		err = AddObjectToZip(zipWritter, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddObjectToZip(zipWritter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	zipFile, err := zipWritter.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipFile, file)
	if err != nil {
		return err
	}

	return nil
}

func matchesExclude(filename, exclude string) bool {
	if strings.HasPrefix(exclude, "*.") {
		ext := strings.TrimPrefix(exclude, "*")
		return strings.HasSuffix(filename, ext)
	}
	return false
}
