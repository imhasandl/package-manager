package main

import (
	"fmt"
	"log"
	"os"
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
	

	return nil
}

func UpdatePackage(configFile string) error {

	return nil
}
