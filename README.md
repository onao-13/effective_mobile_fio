# Сервис ФИО

Ссылка на [тестовое задание](https://docs.yandex.ru/docs/view?url=ya-disk-public%3A%2F%2FWAog00POrvh1MVzh7%2B5Ke3jk3kEt8W%2B1TuJobamPXLUkoitW%2Fsvs3LTji0US2aVNq%2FJ6bpmRyOJonT3VoXnDag%3D%3D&name=%D0%97%D0%B0%D0%B4%D0%B0%D0%BD%D0%B8%D0%B5.%20Junior.pdf)

## OpenAPI
[Путь к файлу](api/fio.yaml)

## Сборка и запуск проекта

```shell
# сборка
docker build . -t fio-service

# запуск
cd deploy; docker compose up
```

## Тестирование
Перед тестированием необходимо поднять тестовое окружение

```shell
cd test; docker compose up
```

## Схема СУБД
[Путь к файлу](db/migration.sql)