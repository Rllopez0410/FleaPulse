# FleaPulse

FleaPulse is a command-line application written in Go for analyzing Escape from Tarkov flea market pricing data and identifying potentially profitable trader flips.

The program analyzes recorded flea market prices, calculates average minimum prices, and compares them against trader sell prices to identify potential profit opportunities.

> **Note:** FleaPulse uses historical pricing data. Results do not guarantee current flea market prices or profits. It is intended as a guide for identifying items that may be worth watching or searching for on the flea market.

---

## Features

- Retrieve Escape from Tarkov item data
- Retrieve historical flea market pricing data
- Calculate average minimum flea market prices
- Compare flea prices against trader sell prices
- Identify potentially profitable trader flips
- Sort profitable items by profit
- Support multiple Tarkov game modes
- Adjustable API request rate
- Interactive command-line interface

---

## Requirements

To build FleaPulse from source, you will need:

- [Go](https://go.dev/) installed on your system
- An internet connection to retrieve Tarkov market data

You can check that Go is installed with:

```bash
go version
```

---

## Installation

Clone the repository:

```bash
git clone https://github.com/Rllopez0410/FleaPulse.git
```

Enter the project directory:

```bash
cd FleaPulse
```

Run FleaPulse:

```bash
go run .
```

Alternatively, build an executable:

```bash
go build
```

### Linux / macOS

Run the compiled executable with:

```bash
./FleaPulse
```

### Windows

Run the compiled executable with:

```powershell
.\FleaPulse.exe
```

---

## Usage

When FleaPulse starts, you will be presented with the command-line interface:

```text
███████╗██╗     ███████╗ █████╗ ██████╗ ██╗   ██╗██╗     ███████╗███████╗
██╔════╝██║     ██╔════╝██╔══██╗██╔══██╗██║   ██║██║     ██╔════╝██╔════╝
█████╗  ██║     █████╗  ███████║██████╔╝██║   ██║██║     ███████╗█████╗
██╔══╝  ██║     ██╔══╝  ██╔══██║██╔═══╝ ██║   ██║██║     ╚════██║██╔══╝
██║     ███████╗███████╗██║  ██║██║     ╚██████╔╝███████╗███████║███████╗
╚═╝     ╚══════╝╚══════╝╚═╝  ╚═╝╚═╝      ╚═════╝ ╚══════╝╚══════╝╚══════╝

Welcome To FleaPulse...

Type [ man ] to see a list of general commands.

Tarkov =>
```

To view the available commands:

```text
Tarkov => man
```

---

## Commands

### `man`

Displays the FleaPulse manual and a list of general commands.

```text
Tarkov => man
```

---

### `pulse`

Scans the flea market for pricing data, calculates the average minimum price of each item, and compares it to the item's trader sell price to identify potential profit opportunities.

```text
Tarkov => pulse
```

Depending on the configured request rate, a full scan may take some time to complete.

---

### `profits`

Retrieves a list of potentially profitable items discovered during a scan.

```text
Tarkov => profits
```

Run `pulse` before using this command so FleaPulse has pricing data to analyze.

---

### `items`

Searches for essential information about each item, such as its ID, name, and type.

```text
Tarkov => items
```

---

### `flea`

Searches for a specified item and returns its pricing history.

```text
Tarkov => flea
```

---

### `mode`

Displays the current game mode and allows you to change it.

```text
Tarkov => mode
```

Available game modes include:

```text
pve
regular
pvp-season
```

---

### `rate`

Displays and allows you to change the API request rate in milliseconds.

```text
Tarkov => rate
```

The request rate controls the delay between requests made while FleaPulse is retrieving pricing data.

For example:

```text
50ms
100ms
250ms
500ms
```

A lower value results in requests being made more frequently, while a higher value increases the delay between requests.

If requests begin failing or the API begins rate limiting requests, increasing the request rate may improve reliability.

---

### `kill`

Terminates FleaPulse.

```text
Tarkov => kill
```

---

## Basic Workflow

A typical FleaPulse session might look like:

```text
Tarkov => mode
Tarkov => rate
Tarkov => pulse
Tarkov => profits
```

1. Select the desired Tarkov game mode.
2. Configure the API request rate if necessary.
3. Run `pulse` to collect and analyze flea market data.
4. Run `profits` to view potential trader flips.

---

## How Profit Is Calculated

FleaPulse records flea market pricing data and calculates the average minimum flea market price for an item.

That value is then compared against the item's trader sell price.

Conceptually:

```text
Potential Profit = Trader Sell Price - Average Minimum Flea Price
```

For example:

```text
Average minimum flea price:  45,000 RUB
Trader sell price:           52,000 RUB
Potential profit:             7,000 RUB
```

In this example, purchasing the item for approximately `45,000 RUB` and selling it to the trader for `52,000 RUB` would represent a potential profit of:

```text
7,000 RUB
```

However, because FleaPulse analyzes recorded pricing data, the item may not currently be available at the calculated flea market price.

---

## Historical Pricing

FleaPulse uses recorded flea market pricing data rather than guaranteeing real-time flea market listings.

Because flea market prices can change, results should be treated as potential opportunities rather than guaranteed trades.

FleaPulse is intended to help identify items that may be worth:

- Watching on the flea market
- Searching for during raids
- Comparing against trader prices
- Investigating for potential profit opportunities

---

## API Request Rate

FleaPulse allows the delay between API requests to be adjusted.

This can be useful if network conditions or API availability change.

For example, a request interval of:

```text
50ms
```

allows requests to occur more frequently than:

```text
500ms
```

If FleaPulse encounters frequent network errors or rate limiting, increasing the delay between requests may help.

Users should avoid placing unnecessary load on the API service.

---

## Data Source

FleaPulse uses data provided by the [tarkov.dev](https://tarkov.dev/) community API.

Market information retrieved through the API may occasionally be delayed, incomplete, or unavailable.

FleaPulse is not affiliated with tarkov.dev, Battlestate Games, or Escape from Tarkov.

---

## Built With

FleaPulse is written primarily in:

- **Go**
- Go standard library HTTP client
- JSON APIs
- Command-line interface tools

The project uses HTTP requests to retrieve Tarkov data and Go structs to decode and analyze the returned JSON.

---

## Project Status

FleaPulse is currently under development.

Features, commands, calculations, output formatting, and API behavior may change as development continues.

---

## Disclaimer

FleaPulse is an unofficial community project.

It is not affiliated with, endorsed by, or associated with Battlestate Games.

Escape from Tarkov and related names, trademarks, and assets are property of their respective owners.

FleaPulse does not guarantee flea market prices, trader profits, or item availability.
