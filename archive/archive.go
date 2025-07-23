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

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

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

		err := AddObjectsToZip(zipWriter, path, exclude)
		if err != nil {
			return err
		}
	}

	fmt.Printf("ZIP архив создан: %s\n", archiveFile)
	return nil
}

func AddObjectsToZip(zipWriter *zip.Writer, path, exclude string) error {
	files, err := filepath.Glob(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if exclude != "" && matchesExclude(file, exclude) {
			continue
		}

		err = AddObjectToZip(zipWriter, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddObjectToZip(zipWriter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	zipFile, err := zipWriter.Create(filename)
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

func ExtractZipArchive(localPath, extractPath string) error {
	zipReader, err := zip.OpenReader(localPath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	err = os.MkdirAll(extractPath, 0755)
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		err = extractFileFromZip(file, extractPath)
		if err != nil {
			return err
		}
	}

	fmt.Printf("ZIP архив распакован в: %s\n", extractPath)
	return nil
}

func extractFileFromZip(file *zip.File, extractPath string) error {
	zipFile, err := file.Open()
	if err != nil {
		return err
	}
	defer zipFile.Close()

	targetPath := filepath.Join(extractPath, file.Name)

	err = os.MkdirAll(filepath.Dir(targetPath), 0755)
	if err != nil {
		return err
	}

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, zipFile)
	if err != nil {
		return err
	}
	return nil
}
