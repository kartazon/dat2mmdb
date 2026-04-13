# dat2mmdb

Converts [v2fly/geoip](https://github.com/v2fly/geoip) `geoip.dat` to a MaxMind-compatible `.mmdb` file suitable for use with [Stash](https://stash.wiki) and other tools that support the **MaxMind GeoIP format**.

## How it works

1. `geoip.dat` (v2ray/v2fly format) is parsed by the official `v2fly/geoip` tool into a text file with CIDR ranges per country tag.
2. `cmd/txt2mmdb` reads that text and writes a `GeoLite2-Country`-compatible `.mmdb` using [`mmdbwriter`](https://github.com/maxmind/mmdbwriter).
3. GitHub Actions runs this pipeline automatically every Monday at 03:00 UTC, or on manual trigger.

## Repository structure

```
.
├── geoip.dat                        # Put your v2fly geoip.dat here
├── go.mod
├── cmd/
│   └── txt2mmdb/
│       └── txt2mmdb.go              # dat → text → mmdb converter
└── .github/
    └── workflows/
        └── build-mmdb.yml           # GitHub Actions pipeline
```

## Usage

### 1. Add geoip.dat

Place your `geoip.dat` file in the root of the repository. You can download the latest one from [v2fly/geoip releases](https://github.com/v2fly/geoip/releases).

### 2. Run locally

```bash
# Install v2fly geoip tool
go install github.com/v2fly/geoip@latest

# Convert dat → text
geoip -c config.json

# Convert text → mmdb
mkdir -p output
go run ./cmd/txt2mmdb -in geoip.txt -out output/geoip-country.mmdb
```

### 3. Run via GitHub Actions

Go to **Actions → build-mmdb → Run workflow**. The resulting `.mmdb` will be available as a workflow artifact.

## Output MMDB schema

Each IP network is stored with a GeoLite2-Country-compatible record:

```json
{
  "country": {
    "iso_code": "CN",
    "names": { "en": "China" }
  },
  "registered_country": {
    "iso_code": "CN",
    "names": { "en": "China" }
  }
}
```

This is compatible with the **MaxMind GeoIP format** expected by Stash.

## License

MIT
