
```
$ tree internal
internal
├── domains
│   ├── availability.go
│   └── order.go
├── infrastructure
│   ├── logger
│   │   └── logger.go
│   └── unitofwork
│       └── mutex_unit_of_work.go
├── presentation
│   ├── httphandlers
│   │   ├── orders_http_handler.go
│   │   └── orders_middlewares.go
│   └── httpmodels
│       └── orders_models.go
├── repositories
│   └── availability_repository.go
├── services
│   └── booking_service.go
└── shared
```

```
$ tree cmd
cmd
└── server
    ├── main.go
    └── main_test.go
```

```
$ go test ./...  -v
```

We can use linux core utils magic to make json slog pretty, see in the end.
```
$ go run cmd/server/main.go
```

Run this twice to validate rejection of overbooking
```
curl -sS --location --request POST 'localhost:8080/orders' \
--header 'Content-Type: application/json' \
--data-raw '{
    "hotel_id": "reddison",
    "room_id": "lux",
    "email": "guest@mail.ru",
    "from": "2024-01-02T00:00:00Z",
    "to": "2024-01-04T00:00:00Z"
}' | jq -R 'fromjson? // .'
```

Run this to test handling required fields by middlewares
```
curl -sS --location --request POST 'localhost:8080/orders' \
--header 'Content-Type: application/json' \
--data-raw '{
    "room_id": "lux",
    "email": "guest@mail.ru",
    "from": "2024-01-02T00:00:00Z",
    "to": "2024-01-04T00:00:00Z"
}' | jq -R 'fromjson? // .'
```

TODO:
- [ ] TODO in the code
- [ ] use order domain, unused for simplicity
- [ ] instead of writing to response jsonBody - create and use OrderAPIResponse struct 
- [ ] various small oversights


Hardcore but very pretty logging, won't work with the Air.
```
go run cmd/server/main.go | unbuffer -p awk -F'"' -v OFS='"' '{ gsub("T"," ",$4); gsub("[0-9]{4}-[0-9]{2}-[0-9]{2} ","",$4); gsub("\\+[0-9:]+","",$4); print }' | jq '.'
```
