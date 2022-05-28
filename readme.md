# Stock Monitor

Simple cli tool to keep track of a portfolio.

## Usage

### Add shares
```
stock-monitor buy <ticker> <price> <number>

Example:
stock-monitor buy MO 27.2351 <number>
```

### Sell shares
```
stock-monitor sell <ticker> <price> <number>

Example:
stock-monitor sell MO 27.2351 <number>
```

### Rename ticker
```
stock-monitor rename <old_ticker> <new_ticker>

Example:
stock-monitor rename CTL LUMN
```

### Show history of actions
```
stock-monitor oh
```

### Show list of current positions
```
stock-monitor s
```

get a token from finnhub.io and set is as FINNHUB_TOKEN environment variable to get current values of positions

## CSV Import

```
stock-monitor import /path/to/file.csv
```

CSV has to be in following format:

```
type,date,ticker,new_ticker,price,number_of_shares

example:
buy,2000-01-01,MO,,12.3456,100
sell,2000-01-01,MO,,12.3456,100
rename,2000-01-01,MO,FOO,,
```
