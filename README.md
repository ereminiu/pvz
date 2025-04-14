## Инструкция

### запуск
```bash
make build
make run
```

### миграция (можно не накатывать при первом запуске, если папка с volume пуста)
```bash
make migrate-up
make migrate-down
```

### UML-схема работы с кэшом
![screenshot](screenshot.png)

### curl запрос

<details>
<summary>Авторизация /sign-in </summary>

```curl
curl --location 'http://localhost:9000/api/sign-in' \
--header 'Content-Type: application/json' \
--data '{
    "username": "Alina Rin",
    "password": "qwerty"
}'
```
</details>


<details>
<summary>Добавить заказ /add </summary>

```curl
curl --location 'http://localhost:9000/api/add' \
--header 'Authorization: Bearer TOKEN_HERE' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": 19,
    "order_id": 12,
    "expire_after": 2,
    "weight": 5,
    "price": 100,
    "packing": "box",
    "extra": true
}'
```
</details>


<details>
<summary>Удаление заказа /remove </summary>

```curl
curl --location --request DELETE 'http://localhost:9000/api/remove' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDEzMzEwMzEsImlhdCI6MTc0MTI4NzgzMSwidXNlcl9yb2xlIjoiYWRtaW4ifQ.zDu67gofZiXusSNqcvcHucAKFFPTz2i0inlQOc6Ku7A' \
--header 'Content-Type: application/json' \
--data '{
    "order_id": 89
}'
```
</details>


<details>
<summary>Вернуть заказы /refund </summary>

```curl
curl --location 'http://localhost:9000/api/refund' \
--header 'Authorization: Bearer TOKEN_HERE' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": 19,
    "order_ids": [89, 86, 84]
}'
```
</details>


<details>
<summary>Получить заказ /return </summary>

```curl
curl --location 'http://localhost:9000/api/return' \
--header 'Authorization: Bearer TOKEN_HERE' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": 19,
    "order_ids": [89, 86, 84]
}'
```
</details>


<details>
<summary>Получение списка заказов по ID пользователя /list </summary>

```curl
curl --location 'http://localhost:9000/api/list' \
--header 'Authorization: Bearer TOKEN_HERE' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": 19,
    "last_n": 5,
    "located": true,
    "pattern": {
        "status": "re",
        "packing": "bo"
    }
}'
```
</details>


<details>
<summary>Получить список возвратов /refund-list </summary>

```curl
curl --location 'http://localhost:9000/api/refund-list' \
--header 'Authorization: Bearer TOKEN_HERE' \
--header 'Content-Type: application/json' \
--data '{
    "page": 3,
    "limit": 10,
    "pattern": {
        "status": "re",
        "packing": "bo"
    }
}'
```
</details>


<details>
<summary>Получить историю заказов /history </summary>

```curl
curl --location 'http://localhost:9000/api/history' \
--header 'Authorization: Bearer TOKEN_HERE' \
--header 'Content-Type: application/json' \
--data '{
    "page": 1,
    "limit": 3,
    "pattern": {
        "status": "re",
        "packing": "bo"
    }
}'
```
</details>

## Домашнее задание
Работу продолжаем в репозитории Домашнего Задания 2.

## Сохранение данных в БД
Подключить работу с базой данных, добавить конфиг подключения, инициализацию коннекта.

БД должна быть реализована в 3 нормальной форме.

Добавить фильтрацию на все list ручки - фильтры должны конфликтовать (к примеру, зона доставки)

Написать CRUD операции для работы с БД.
Должны быть реализованы методы записи и чтения данных простой системы хранения ПВЗ.

Подсказка:

Вам могут пригодиться следующие методы

- GetByID
- List
- Update
- Create
- Delete

Также необходимо реализовать те методы, которые вы делали с файлом.

## Разработать HTTP сервер.

Необходимо реализовать HTTP сервер, который будет работать с методами базы данных, реализованными в 1 пункте.

Методы должны позволять манипулировать данными (create, read, update, delete) для системы хранения ПВЗ.

Методы должны быть выполнены в restful стиле. Необходимо корректно обрабатывать все коды ошибок.

Входящие параметры должны быть представлены либо в формате json, либо в query параметрах (зависит от метода).

Сервис должен использовать порт 9000.

В ридми приложить curl запросы на каждую ручку. Запросы должны быть валидными и возвращать нужный код ответа.
Реализовать middleware с basic auth. Юзер/пароль можно задать как в конфиге, так и хранить в базе (создать круд юзеров).

## Подсказка:

Посмотрите на результат выполнения ДЗ 2. Сервис должен делать тот же flow, но используя БД как хранилище и http как интерфейс взаимодействия с пользователем.

## Дополнительное задание
- Необходимо реализовать middleware, который будет логгировать поля POST,PUT,DELETE запросов
- На методы "Получения списка..." реализовать пагинацию


Ограничения дз:

- Нельзя использовать orm или sql билдеры
- Для реализации http сервера можно использовать как net/http, так и gin/fasthttp и прочее
- Коды ошибок должны соответствовать поведению сервиса. Хендлеры, которые отдают только 500 в случае ошибки - не принимаются
- В хендлерах должна быть базовая валидация данных, соответствующая бизнес-логике
- Пагинация должна работать курсором

### Дедлайны сдачи и проверки задания: 
- 8 марта 23:59 (сдача) / 11 марта, 23:59 (проверка)

