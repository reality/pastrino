# Pastrino

You can use this application to show recent news stories related to stocks you own or are interested in. It will first construct your portfolio from the options you pass it, then list them with prices from the Yahoo finance API. Then it will search through a list of RSS feeds for mentions of the stocks, then listing them associated with a binary pos/neg sentiment score.

## Usage

```bash
go run main.go -t212 [filename] -wl [filename]
```

### Arguments

#### -t212

Here you can pass the CSV export of transactions from your T212 account (you can get it from the transaction history page), It will construct the current state of your ownership by 'simulating' your history of buys and sells.

#### -wl

Here you can pass the path to a file that contains a manually constructed list of stocks, with associated keywords to search for in news items. You can also extend stock entries produced by other methods with additional keywords.

The format of the file should be

```tsv
Ticker	Name	Keywords
VFF	Village Farms International	cannabis,biscuits
```

You can see an example in the `example_watchlist.tsv` file.

## Configuration

There is also a `config.json` file, in which you can modify the list of RSS feeds to scan for news.
