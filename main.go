package main

import (
	"fmt"
	"log"
	"os"

	"github.com/imhasandl/package-manager/archive"
	"github.com/imhasandl/package-manager/config"
)

func main() {
	fmt.Println("Использование:")
	fmt.Println("pm create <packet.json>")
	fmt.Println("pm update <packages.json>")

	command := os.Args[1]

	switch command {
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("pm create <package.json>")
		}
		err := CreatePackage(os.Args[2])
		if err != nil {
			log.Printf("Не получилось создать архив: %v", err)
		}
	case "update":
		if len(os.Args) < 3 {
			fmt.Println("pm update <package.json>")
		}
		err := UpdatePackage(os.Args[2])
		if err != nil {
			log.Printf("Не получилось обновить архив: %v", err)
		}
	}
}

func CreatePackage(configFile string) error {
	config, err := config.ReadPackageConfig(configFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения конфига:", err)
	}

	packageName := config["name"].(string)
	packageVersion := config["ver"].(string)
	archiveFile := fmt.Sprintf("%s-%s.zip", packageName, packageVersion)

	err = archive.CreateZipArchive(config, archiveFile)
	if err != nil {
		return fmt.Errorf("ошибка создания архива: %v", err)
	}

	return nil
}

func UpdatePackage(configFile string) error {

	return nil
}
