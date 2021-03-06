# Tracker dev server

В этом репозитории находится dev сервер для микросервиса tracker

## Локальаня разработка

### Необходимые инструменты

Для разработки микросервиса нужен локально установленный GO версии `1.17+`

### Настройка переменных GO

```bash
go env -w GO111MODULE=on

```

### Enviroments

В корне репозитория находится файл `.env.example`. С его содержимым нужно создать файл в корне репозиторя `.env` и заменить нужные переменные на свои

#### Описание переменных

| Name                | Description                       | Default value | Available values                                      |
| ------------------- | --------------------------------- | ------------- | ----------------------------------------------------- |
| `APP_NAME`          | Полное название микросервиса      | `Tracker`     | `string`                                              |
| `APP_MS_NAME`       | Сокращенное Название микросервиса | `tracker-dev` | `string`                                              |
| `APP_HOST`          | Хост сервера                      | `localhost`   | `string`                                              |
| `ENV`               | Enviroment                        | `dev`         | `dev`, `stage`, `prod`                                |
| `LOG_LEVEL`         | Уровень логирования               | `info`        | `debug`, `info`, `warning`, `error`, `fatal`, `panic` |
| `GRPC_PORT`         | GRPC port                         | `9001`        | `string`                                              |
| `PROXY_DEV_KEY`     | Dev key for proxy server          |               | `string`                                              |
| `PROXY_REMOTE_HOST` | Relay server hostname             |               | `string`                                              |

### Запуск микросервиса

```bash
make start
```

### Запуск проксирующего сервера

```bash
make proxyclient port=35100
```

> Где `port` - порт нужного протокола

Принцип работы проксирующего сервера следующий:

После того как к нему подключается текущий микросервис, на проксируемом открывается порт для прослушки клиентов (трекеров).
Далее, после установки соединения клиента с проксируемым сервером, все данные без каких либо обработок будут пересылаться на текущий микросервис.

Клиентов может быть любое количество, проксироваться данные могут только на один удаленный сервер в рамках одного протокола одновременно. Это сделано для того, чтобы разработчик получал все данные со всех трекеров и только он мог отвечать клиентам

На проксируемом сервере открыты порты для всех протоколов

> Для авторизации на проксирующем сервере нужны переменные окружения `PROXY_DEV_KEY` и `PROXY_REMOTE_HOST` в файле `.env`
