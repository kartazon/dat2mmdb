# dat2mmdb

Converts `geoip.dat` (v2fly/v2ray format) from [Ground-Zerro/Geo-Aggregator](https://github.com/Ground-Zerro/Geo-Aggregator) to a MaxMind-compatible `.mmdb` file suitable for use with [Stash](https://stash.wiki) and other tools that support the **MaxMind GeoIP format**.

## How it works

1. `geoip.dat` is automatically downloaded from [Ground-Zerro/Geo-Aggregator](https://raw.githubusercontent.com/Ground-Zerro/Geo-Aggregator/main/geodat/geoip_GA.dat) — no manual upload needed.
2. The file is parsed by the official [v2fly/geoip](https://github.com/v2fly/geoip) tool into a text file with CIDR ranges per country tag.
3. `cmd/txt2mmdb` reads that text and writes a `GeoLite2-Country`-compatible `.mmdb` using [`mmdbwriter`](https://github.com/maxmind/mmdbwriter).
4. GitHub Actions runs this pipeline automatically every Monday at 03:00 UTC, or on manual trigger.

## Repository structure

```
.
├── go.mod
├── cmd/
│   └── txt2mmdb/
│       └── txt2mmdb.go              # dat → text → mmdb converter
└── .github/
    └── workflows/
        └── build-mmdb.yml           # GitHub Actions pipeline
```

## Usage

### Run via GitHub Actions

Go to **Actions → build-mmdb → Run workflow**.

`geoip.dat` is fetched automatically from:
```
https://raw.githubusercontent.com/Ground-Zerro/Geo-Aggregator/main/geodat/geoip_GA.dat
```

The resulting `geoip-country.mmdb` will be available as a workflow artifact.

### Run locally

```bash
# Download geoip.dat
curl -fsSL https://raw.githubusercontent.com/Ground-Zerro/Geo-Aggregator/main/geodat/geoip_GA.dat -o geoip.dat

# Install v2fly geoip tool
go install github.com/v2fly/geoip@latest

# Convert dat → text
geoip -c config.json

# Convert text → mmdb
mkdir -p output
go run ./cmd/txt2mmdb -in geoip.txt -out output/geoip-country.mmdb
```

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
