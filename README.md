# osmtools

Spreadsheet conversion tools for OSM accountancy import. All wrapped up in a web app.

## Overview

`osmtools` is a small Go web application that cleans up CO-OP Bank statement
spreadsheets (`.xlsx`) so they can be imported into the OSM accounting system.
It serves a simple web UI where a statement file can be uploaded, processed,
and downloaded again in OSM-compatible format.

## Features

- **30 day statement conversion** (`POST /coop30day-osm`) — strips the fixed
  header/footer rows and unused columns from a standard CO-OP 30 day
  statement export.
- **Custom statement conversion** (`POST /coopcustom-osm`) — same cleanup for
  the CO-OP custom statement layout, which has a slightly different row/column
  structure.
- Both conversions remove unneeded rows/columns, resize columns, and strip an
  embedded logo picture from the spreadsheet before returning the cleaned
  file for download.
- A single HTML page (`templates/index.html`, styled with
  [Bulma](https://bulma.io/)) provides file upload forms for both conversion
  types.

## Project layout

| Path                    | Description                                                        |
|-------------------------|---------------------------------------------------------------------|
| `main.go`                | HTTP server setup and route registration                          |
| `handlers.go`            | Request handlers and spreadsheet processing logic (uses [excelize](https://github.com/xuri/excelize)) |
| `templates/index.html`   | Web UI for uploading statements                                    |
| `Dockerfile`             | Multi-stage build for a minimal distroless container image        |
| `.github/workflows`      | CI/CD workflow definitions                                          |

## Requirements

- Go 1.24+ (see `go.mod`)

## Running locally

```powershell
go run .
```

The server listens on port `8000`. Open `http://localhost:8000` in a browser
to use the upload forms.

## Building

```powershell
go build .
```

## Running with Docker

```powershell
docker build -t osmtools .
docker run -p 8000:8000 osmtools
```
