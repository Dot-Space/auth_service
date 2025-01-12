# Project Structure

The project is organized as follows:
```
непостоянные конфиги:
    - /config/config.*.yaml
    - .env
    ( добавленны в .gitignore )

корневые конфиги:
    - .gitignore
    - docker-compose.*.yml
    - go.mod
    - go.sum
    - Makefile
    - Dockerfile
    - .air.toml

/ cmd
    - подготовка, инициализация запуска
    - миграции
    - подключение логера

/ config
    - окружение
    - конфиги сервера
    ( есть не постоянные конфиги, config.*.yaml )

/ gen-proto
    - генерация gRPC сервера
    - cхема protoBuf

/ internal
    = внутренняя логика сервиса =
    / db
        - подключение к базе данных
        - sql CRUD's

    / models
        - структуры данных

    / grpc
        - реализация gRPC методов
        - хэндлеры

    / usecases
        - бизнес логика

    /
```