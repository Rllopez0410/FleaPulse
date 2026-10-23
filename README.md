# FleaPulse

**Current Version: v1.0.0**

FleaPulse is a command-line application written in Go for analyzing Escape from Tarkov flea market pricing data and identifying potentially profitable trader flips.

The program analyzes recorded flea market prices, calculates average minimum prices, and compares them against trader sell prices to identify potential profit opportunities.

> **Note:** FleaPulse uses historical pricing data. Results do not guarantee current flea market prices or profits. It is intended as a guide for identifying items that may be worth watching or searching for on the flea market.

---

## Download

FleaPulse v1.0.0 is available as a precompiled executable.

You **do not need Go installed** to use a precompiled version.

Go to the **Releases** section of this repository and download the appropriate version for your operating system:

| Operating System | Download |
| --- | --- |
| Windows 64-bit | `FleaPulse-windows-amd64.exe` |
| Linux 64-bit | `FleaPulse-linux-amd64` |
| macOS Intel 64-bit | `FleaPulse-macos-amd64` |

> **macOS:** The current macOS release is built for Intel (`amd64`) Macs.

---

## Running FleaPulse

### Windows

Download:

```text
FleaPulse-windows-amd64.exe
```

Open a terminal in the directory containing the executable and run:

```powershell
.\FleaPulse-windows-amd64.exe
```

### Linux

Download:

```text
FleaPulse-linux-amd64
```

You may need to make the file executable:

```bash
chmod +x FleaPulse-linux-amd64
```

Then run:

```bash
./FleaPulse-linux-amd64
```

### macOS — Intel

Download:

```text
FleaPulse-macos-amd64
```

You may need to make the file executable:

```bash
chmod +x FleaPulse-macos-amd64
```

Then run:

```bash
./FleaPulse-macos-amd64
```

> Your operating system may display a security warning because FleaPulse binaries are currently unsigned.

---

## Features

- Retrieve Escape from Tarkov item data
- Retrieve recorded flea market pricing data
- Calculate average minimum flea market prices
- Compare flea prices against trader sell prices
- Identify potentially profitable trader flips
- Sort profitable items by profit
- Support multiple Tarkov game modes
- Adjustable API request rate
- Interactive command-line interface

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

FleaPulse v1.0.0

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

### `pulse`

Scans the flea market for pricing data, calculates the average minimum price of each item, and compares it to the item's trader sell price to identify potential profit opportunities.

```text
Tarkov => pulse
```

Depending on the configured request rate, a full scan may take some time to complete.

### `profits`

Retrieves a list of potentially profitable items discovered during a scan.

```text
Tarkov => profits
```

Run `pulse` before using this command so FleaPulse has pricing data to analyze.

### `items`

Searches for essential information about each item, such as its ID, name, and type.

```text
Tarkov => items
```

### `flea`

Searches for a specified item and returns its pricing history.

```text
Tarkov => flea
```

### `mode`

Displays the current game mode and allows you to change it.

```text
Tarkov => mode
```

Available game modes:

```text
pve
regular
pvp-season
```

### `rate`

Displays and allows you to change the API request rate in milliseconds.

```text
Tarkov => rate
```

The request rate controls the delay between requests made while FleaPulse retrieves pricing data.

For example:

```text
50ms
100ms
250ms
500ms
1000ms
```

A lower value results in requests being made more frequently, while a higher value increases the delay between requests.

If FleaPulse encounters frequent network errors or API rate limiting, increasing the delay may improve reliability.

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

That value is compared against the item's trader sell price.

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

In this example, purchasing the item for approximately `45,000 RUB` and selling it to the trader for `52,000 RUB` represents a potential profit of `7,000 RUB`.

Because FleaPulse analyzes recorded pricing data, the item may not currently be available at the calculated flea market price.

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

FleaPulse allows you to adjust the delay between API requests.

This can be useful if network conditions or API availability change.

For example:

```text
50ms
```

makes requests more frequently than:

```text
500ms
```

If FleaPulse encounters frequent network errors or rate limiting, increasing the delay between requests may help.

Please avoid placing unnecessary load on the API service.

---

## Building From Source

If you prefer to build FleaPulse yourself, you will need Go installed.

Check your Go installation:

```bash
go version
```

Clone the repository:

```bash
git clone https://github.com/Rllopez0410/FleaPulse.git
```

Enter the project directory:

```bash
cd FleaPulse
```

Run directly from source:

```bash
go run .
```

Or build the application:

```bash
go build
```

---

## Data Source

FleaPulse uses data provided by the [tarkov.dev](https://tarkov.dev/) community API.

Market information retrieved through the API may occasionally be delayed, incomplete, or unavailable.

FleaPulse is not affiliated with tarkov.dev, Battlestate Games, or Escape from Tarkov.

---

## Built With

FleaPulse is written in **Go** and makes use of:

- Go standard library HTTP functionality
- JSON APIs
- Go structs and collections
- Command-line interface tools

---

## Project Status

**FleaPulse v1.0.0 is the first stable release of the project.**

Core development is complete, but FleaPulse will continue to receive new features, improvements, optimizations, and bug fixes over time.

---

## Release History

### v1.0.0

First stable release of FleaPulse.

Highlights include:

- Flea market price scanning
- Historical price analysis
- Trader profit calculations
- Profitable item discovery
- Multiple game mode support
- Configurable API request rate
- Interactive CLI and built-in manual
- Precompiled Windows, Linux, and Intel macOS releases

---

## Disclaimer

FleaPulse is an unofficial community project.

It is not affiliated with, endorsed by, or associated with Battlestate Games.

Escape from Tarkov and related names, trademarks, and assets are property of their respective owners.

FleaPulse does not guarantee flea market prices, trader profits, or item availability.
