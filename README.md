# WeatherTracker

Простой бэкенд-сервис на Go для отслеживания погодных данных.

Проект периодически получает данные о температуре из открытого API [Open-Meteo](https://open-meteo.com/), сохраняет их в PostgreSQL и предоставляет доступ через HTTP-сервер.

## Особенности

- Автоматический сбор данных о погоде по расписанию (с использованием gocron)
- Хранение данных в PostgreSQL
- HTTP-сервер на базе chi
- Клиент для работы с Open-Meteo API
- Запуск базы данных через Docker Compose

## Технологии

- **Язык**: Go (1.25+)
- **Веб-фреймворк**: [chi](https://github.com/go-chi/chi)
- **Планировщик задач**: [gocron](https://github.com/go-co-op/gocron)
- **Драйвер БД**: [pgx](https://github.com/jackc/pgx)
- **База данных**: PostgreSQL
- **Контейнеризация**: Docker Compose

## Установка и запуск

1. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/drenk83/WeatherTracker.git
   cd WeatherTracker
   ```

2. Запустите базу данных:

   ```bash
   docker-compose up -d
   ```

   Это поднимет контейнер PostgreSQL с параметрами:
   - Порт: `54321` (маппится на 5432 внутри контейнера)
   - Пользователь: `drenk83`
   - Пароль: `password`
   - База данных: `weather`

3. Создайте таблицу для хранения данных:

   Подключитесь к БД (например, через psql):

   ```bash
   psql -h localhost -p 54321 -U drenk83 -d weather
   ```

   Пароль: `password`

   Выполните SQL:

   ```sql
   CREATE TABLE reading (
       name TEXT NOT NULL,
       time TIMESTAMP NOT NULL,
       temperature FLOAT8 NOT NULL
   );
   ```

4. Запустите приложение:

   ```bash
   go run ./cmd/server
   ```

   Приложение запустит HTTP-сервер и планировщик задач для периодического получения данных о погоде.