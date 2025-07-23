# Package Manager

Простой менеджер пакетов на Go для создания ZIP-архивов и управления пакетами через SSH.

## Что делает

- Создает ZIP-архивы из файлов и загружает их на сервер
- Скачивает и устанавливает пакеты с сервера
- Управляет зависимостями через JSON конфигурацию

## Как запустить

### Установка
```bash
git clone <repository-url>
cd package-manager
go mod tidy
go build -o pm main.go
```

### Использование
```bash
# Создать и загрузить пакет
pm create packet.json

# Установить пакеты
pm update packages.json
```

## Примеры конфигурации

### packet.json (для создания пакета)
```json
{
  "name": "my-package",
  "ver": "1.0.0",
  "targets": [
    "src/*.go",
    "README.md"
  ]
}
```

### packages.json (для установки пакетов)
```json
{
  "packages": [
    {
      "name": "package1",
      "ver": "1.0.0"
    },
    {
      "name": "package2"
    }
  ]
}
```

## Настройка SSH

Отредактируйте `ssh/ssh.go` с вашими данными для подключения к серверу.

Пакеты сохраняются в папку `./packages/`
