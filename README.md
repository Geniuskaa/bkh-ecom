## Тестовое задание БКХ Еком.

Сервер с API можно запустить командой 
```shell
docker compose up
```

- GET /api/counter/:bannerID
- POST /api/stats/:bannerID

Пример тела запроса `POST /api/stats/20`:

```json
{
  "tsFrom": "2024-01-01T00:00:00Z",
  "tsTo": "2024-12-31T23:59:59Z"
}
```