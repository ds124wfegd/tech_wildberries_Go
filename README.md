## Разработан демонстрационный сервис (Go, Postgres, Kafka) с простейшим интерфейсом, отображающий данные о заказе.

### Сервис получает данные заказов из очереди (Kafka), сохраняет их в базу данных (PostgreSQL) и кэширует в памяти для быстрого доступа.

Использованы следующие пакеты/технологии:
Технология | Применение
:----------|:--------:|
Visual Studio Code | редактор кода
PostgreSQL | в качестве хранилища данных
Docker-compose |для запуска образа PostgreSQL, Kafka, zookeeper
Postman |тестирование сценариев API
github.com/spf13/viper| чтение из конфига
github.com/jmoiron/sqlx|для настройки PostgreSQL
github.com/jackc/pgx/v5 v5.7.5|драйвер PostgreSQL
github.com/sirupsen/logrus|для логирования
github.com/gin-gonic/gin| высокопроизводительный веб-фреймворк HTTP (для маршрутизации)
github.com/segmentio/kafka-go v0.4.48|пакет для работы с Kafka

# Решение
* Сервис REST API написан на Golang v 1.23.12 с использованием Clean Architecture, что позволяет легко расширять функционал сервиса и тестировать его. Также был реализован Graceful Shutdown для корректного завершения работы сервиса.

* Реализовано кэширование данных в сервисе: хранятся последние полученные данные о 100 заказах в map, чтобы быстро выдавать их по запросу.

* При перезапуске восстанавливается кеш из БД

* Реализован HTTP-эндпоинт:
```
http://localhost:8080/order/order?order_uid=<order_id>
```
который по order_id возвращает данные заказа из кеша, если в кеше данных нет, то подтягивает из БД.

* Разработан простой веб-интерфейс — страница (HTML/JS), где можно ввести ID заказа и получить информацию о нём, обращаясь к вышеописанному HTTP API, доступна по адресу:
```
http://localhost:8080/
```
* Модель данных заказа:
```
{
   "order_uid": "b563feb7b2b84b6test",
   "track_number": "WBILMTESTTRACK",
   "entry": "WBIL",
   "delivery": {
      "name": "Test Testov",
      "phone": "+9720000000",
      "zip": "2639809",
      "city": "Kiryat Mozkin",
      "address": "Ploshad Mira 15",
      "region": "Kraiot",
      "email": "test@gmail.com"
   },
   "payment": {
      "transaction": "b563feb7b2b84b6test",
      "request_id": "",
      "currency": "USD",
      "provider": "wbpay",
      "amount": 1817,
      "payment_dt": 1637907727,
      "bank": "alpha",
      "delivery_cost": 1500,
      "goods_total": 317,
      "custom_fee": 0
   },
   "items": [
      {
         "chrt_id": 9934930,
         "track_number": "WBILMTESTTRACK",
         "price": 453,
         "rid": "ab4219087a764ae0btest",
         "name": "Mascaras",
         "sale": 30,
         "size": "0",
         "total_price": 317,
         "nm_id": 2389212,
         "brand": "Vivienne Sabo",
         "status": 202
      }
   ],
   "locale": "en",
   "internal_signature": "",
   "customer_id": "test",
   "delivery_service": "meest",
   "shardkey": "9",
   "sm_id": 99,
   "date_created": "2021-11-26T06:22:19Z",
   "oof_shard": "1"
}
```

* Т.к. данные, приходящие из очереди, могут быть невалидными, предусмотрена обработка ошибок.
  
* При различных сбоях (ошибка базы, падение сервиса) данные не теряются — использованы транзакции, механизм подтверждения сообщений от брокера.

* Написан скрипт-эмулятор отправки сообщений в Kafka: на примере убедились, что сервис подключается к брокеру сообщений и обрабатывает сообщения онлайн, для запуска используется команда:
  
```
go run ./scr/producer.go
```

* Кеш действительно ускоряет получение данных (например, при повторных запросах по одному и тому же ID). Было подтверждено с помощью Postman:

* ### без кэша - ~30 мс
  <img width="1374" height="712" alt="image" src="https://github.com/user-attachments/assets/059508c5-094c-410c-bb88-e8fa7a1db947" />
  
* ### с кэшем ~ 10 мс
  <img width="1322" height="713" alt="image" src="https://github.com/user-attachments/assets/03c82c03-bfdf-4ddc-a60a-34295b97e5da" />

* HTTP-сервер возвращает корректные данные в формате JSON.

* Используется Makefile:

Название | Команда | Описание 
:----------|:--------:|:--------:|
build |go build -o bin/order-service ./cmd/app| Сборка приложения
run | go run ./cmd/app |  Запуск приложения
docker-build | docker build -t app . | Сборка Docker образа
up | docker-compose up -d | Запуск PostgreSQL, Kafka
down | docker-compose down | Останов PostgreSQL, Kafka
status| docker-compose ps | Проверка docker status
logs | docker-compose logs -f | Показать логи docker-compose



