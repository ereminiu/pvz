## Инструкция

### запуск
```bash
make build
make run
```

### миграция (можно не накатывать при первом запуске, если нет БД на заданном порту)
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
