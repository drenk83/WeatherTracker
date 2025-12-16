код не очень, согл
подключение к таблице через psql: psql -h localhost -p 54321 -U drenk83 -d weather
создание таблицы: create table reading (name text not null, time timestamp not null, temperature float8 not null); 