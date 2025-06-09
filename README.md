# go_payment_microservice
The payment microservice, implemented by Golang.


### Project setup steps:

1. Open directory with `docker-compose` and `.env` file

2. Open `.env` file and change values

3. Write commands:
   check builders and use desktop-linux

```shell
docker buildx use eloquent_hertz
```

3.1 Build base dockerfile:

```shell
docker build --platform linux/arm64 -f Dockerfile -t base_dockerfile .
```
run directly container without docker compose
```shell
docker run -d --name go_payment -p 8080:8080 base_dockerfile
```

3.2 Build containers:

```shell
docker-compose -f docker-compose.yml build
```

3.3 Run project:

```shell
docker-compose -f docker-compose.yml up -d
```

3.4 Stop project:

```shell
docker-compose -f docker-compose.yml down
```

Generate SWAGGER documentations:

```bash
cd cmd/server && swag init --parseDependency --parseInternal
```

test body for create subscribe
{
"amount": 1,
"card": "5168745017514750",
"card_cvv": "026",
"card_exp_month": "01",
"card_exp_year": "24",
"description": "test",
"order_id": "customer2:1",
"phone": "380936007646",
"subscribe": "1",
"subscribe_date_start": "2025-06-08 00:00:00",
"subscribe_periodicity": "day"
}