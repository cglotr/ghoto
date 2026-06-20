# Ghoto
Uploads and deletes photos to Google Photos

https://cglotr.github.io/ghoto/

[![codecov](https://codecov.io/gh/cglotr/ghoto/graph/badge.svg?token=htlQz2dPW3)](https://codecov.io/gh/cglotr/ghoto)

## What it does

Scans a local directory for `.jpg` and `.mp4` files, uploads them to a Google Photos album, then removes the local files once upload is confirmed. Non-photo files (`.dng`, `.lrv`) are cleaned up automatically. Uploads run in parallel across up to 10 workers, with automatic retry on failure (up to 10 attempts).

## Requirements

- Go 1.25+
- A `credentials.json` file from [Google Cloud Console](https://console.cloud.google.com/) with the Photos Library API enabled

## Setup

1. Create a project in Google Cloud Console and enable the **Photos Library API**
2. Create OAuth 2.0 credentials (Desktop app) and download as `credentials.json` to the project root
3. Run the app — it will open a browser for Google sign-in on first use (OAuth flow via `localhost:8080`)

## Usage

```sh
go run main.go -dir /path/to/photos -album "My Album"
```

Flags:

| Flag | Default | Description |
|------|---------|-------------|
| `-dir` | _(required)_ | Directory containing photos/videos to upload |
| `-album` | `Insta360` | Google Photos album name to upload into |
| `-dryrun` | `false` | Simulate the run without making any changes |

## Dry run

```sh
go run main.go -dir ./testfile/ -dryrun
```

## Testing

```sh
./go__test.sh
```
