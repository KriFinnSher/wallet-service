# **Сервис для управления кошельком**

## Инструкция по запуску
1. Первым делом клонируйте данный репозиторий командой ниже или скачайте архив с исходным кодом и распакуйте:
```bash
git clone https://github.com/KriFinnSher/wallet-service.git
```
2. Далее перейдите в корневую директорию проекта и запустите сервис с помощью [Docker](https://docs.docker.com/get-started/get-docker/):
```bash
cd wallet-service
docker-compose up --build
```
3. Готово! Теперь основной сервер доступен локально на порту `:8080`. Миграции запускаются автоматически при запуске сервиса.

## Тестирование
Для тестирования апи можете использовать curl (в миграциях создаются 2 тестовых пользователя с балансом 1000):
для post-запроса
```bash
curl -X POST http://localhost:8080/api/v1/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "walletId": "00000000-0000-0000-0000-000000000001",
    "operationType": "DEPOSIT",
    "amount": 100
  }'
```
для get-запроса
```bash
curl -X GET http://localhost:8080/api/v1/wallets/00000000-0000-0000-0000-000000000001
```

Также в папке `./tests/load/res` представлены результаты нагрузочного тестирования сервиса.
Конфигурация теста:
- 10_000 запросов (пополам между двумя ручками)
- время теста 30 секунд
<img width="1641" height="619" alt="image" src="https://github.com/user-attachments/assets/372c3c8c-0ec9-47f5-8c58-aae496a7d96f" />



