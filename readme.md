# Stock Monitor

Simple web tool to keep track of a portfolio.

## Usage

- copy env.example to .env and add a finnhub api token
- run `docker-compose up -d`

### Add shares
`POST`

`http://localhost/add-stock`

json payload:
```
{
    "ticker": "FOO",
    "shares": 100,
    "price": 19.99,
    "date": "2023-01-01"
}
```

### Sell shares
`POST`

`http://localhost/sell-stock`

json payload:
```
{
    "ticker": "FOO",
    "shares": 100,
    "price": 19.99,
    "date": "2023-01-01"
}
```

### Rename ticker
`POST`

`http://localhost/rename-stock`

json payload:
```
{
    "old_ticker": "FOO",
    "new_ticker": "BAR",
    "date": "2023-01-01"
}
```

### Show history of orders
`GET`

`http://localhost/order-history`

### Show portfolio
`GET`

`http://localhost/portfolio`

### Show dividends
`GET`

`http://localhost/dividend-history`

Filter by year and/or ticker:

`?year=2023&ticker=FOO`
