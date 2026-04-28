# BCV Scraper

## Prerequisites
Before running the scraper, ensure the following environment variables are set:

* `TELEGRAM_TOKEN`: Your bot's API token.
* `TELEGRAM_CHAT_ID`: The target chat or channel ID.
* `DB_STRING`: Your database connection string (PostgreSQL).

## Usage
To run the application in development mode:
```bash
go run .
```
## Build
To compile the system and generate an executable:
```bash
go build -o bcvScraper . .
```