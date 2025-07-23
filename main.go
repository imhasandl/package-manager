package main

import (
	"fmt"
	"log"
	"os"

	"github.com/imhasandl/package-manager/archive"
	"github.com/imhasandl/package-manager/config"
	"github.com/imhasandl/package-manager/ssh"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование:")
		fmt.Println("  pm create <packet.json>")
		fmt.Println("  pm update <packages.json>")
		return
	}

	command := os.Args[1]

	switch command {
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Использование: pm create <packet.json>")
			return
		}
		err := CreatePackage(os.Args[2])
		if err != nil {
			log.Printf("Не получилось создать архив: %v", err)
		}
	case "update":
		if len(os.Args) < 3 {
			fmt.Println("Использование: pm update <packages.json>")
			return
		}
		err := UpdatePackage(os.Args[2])
		if err != nil {
			log.Printf("Не получилось обновить архив: %v", err)
		}
	default:
		fmt.Printf("Неизвестная команда: %s\n", command)
		fmt.Println("Доступные команды: create, update")
	}
}

func CreatePackage(configFile string) error {
	config, err := config.ReadPackageConfig(configFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения конфига: %v", err)
	}

	packageName := config["name"].(string)
	packageVersion := config["ver"].(string)
	archiveFile := fmt.Sprintf("%s-%s.zip", packageName, packageVersion)

	err = archive.CreateZipArchive(config, archiveFile)
	if err != nil {
		return fmt.Errorf("ошибка создания архива: %v", err)
	}

	remotePath := "/packages/" + archiveFile
	err = ssh.UploadFile(archiveFile, remotePath)
	if err != nil {
		return fmt.Errorf("ошибка загрузки: %v", err)
	}

	fmt.Printf("✓ Пакет %s успешно создан и загружен!\n", packageName)
	return nil
}

func UpdatePackage(configFile string) error {
	config, err := config.ReadPackagesConfig(configFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения конфига: %v", err)
	}

	packages := config["packages"].([]interface{})

	for _, pkg := range packages {
		pkgMap := pkg.(map[string]interface{})
		packageName := pkgMap["name"].(string)

		var packageVersion string
		if ver, exists := pkgMap["ver"]; exists {
			packageVersion = ver.(string)
		}

		err := downloadAndInstallPackage(packageName, packageVersion)
		if err != nil {
			log.Printf("Ошибка установки пакета %s: %v", packageName, err)
			continue
		}
		fmt.Printf("✓ Пакет %s успешно установлен\n", packageName)
	}

	fmt.Println("Обновление завершено!")
	return nil
}

func downloadAndInstallPackage(packageName, packageVersion string) error {
	var archiveFile string
	if packageVersion != "" {
		archiveFile = fmt.Sprintf("%s-%s.zip", packageName, packageVersion)
	} else {
		archiveFile = packageName + "-latest.zip"
	}

	remotePath := "/packages/" + archiveFile
	localPath := "./downloads/" + archiveFile

	err := ssh.DownloadFile(remotePath, localPath)
	if err != nil {
		return fmt.Errorf("не получилось скачать файл: %v", err)
	}

	extractPath := "./packages/" + packageName
	err = archive.ExtractZipArchive(localPath, extractPath)
	if err != nil {
		return fmt.Errorf("не удалось распаковать: %v", err)
	}

	return nil
}