
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
│   ├── apierrors
│   │   └── apierrors.go
│   ├── apihandlers
│   │   ├── order_handlers.go
│   │   └── order_middlewares.go
│   └── apimodels
│       └── order_models.go
├── repositories
│   └── roomavailability
│       └── repository.go
└── services
    └── booking_service.go
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

Run this twice to validate rejection of overbooking.
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

Required fields are enforced, this request will result in 422 error.
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

UTC time is enforced, this request will result in 422 error.
```
curl -sS --location --request POST 'localhost:8080/orders' \
--header 'Content-Type: application/json' \
--data-raw '{
    "hotel_id": "reddison",
    "room_id": "lux",
    "email": "guest@mail.ru",
    "from": "2024-01-01T19:00:00-05:00",
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
