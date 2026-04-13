# dat2mmdb

[English](#english) | [Русский](#русский)

---

## English

Converts `geoip.dat` (v2fly/v2ray format) from [Ground-Zerro/Geo-Aggregator](https://github.com/Ground-Zerro/Geo-Aggregator) to a MaxMind-compatible `.mmdb` file suitable for use with [Stash](https://stash.wiki) and other tools that support the **MaxMind GeoIP format**.

### How it works

1. `geoip.dat` is automatically downloaded from [Ground-Zerro/Geo-Aggregator](https://raw.githubusercontent.com/Ground-Zerro/Geo-Aggregator/main/geodat/geoip_GA.dat) — no manual upload needed.
2. The file is parsed by the official [v2fly/geoip](https://github.com/v2fly/geoip) tool into per-tag text files with CIDR ranges.
3. `cmd/txt2mmdb` reads those files and writes a `GeoLite2-Country`-compatible `.mmdb` using [`mmdbwriter`](https://github.com/maxmind/mmdbwriter).
4. GitHub Actions runs this pipeline automatically every 12 hours, or on manual trigger.

All tags from `geoip.dat` are exported — both country codes (e.g. `CN`, `RU`, `US`) and custom categories (e.g. `NETFLIX`). This allows using arbitrary tags in Stash `GEOIP` rules:

```
GEOIP,CN,DIRECT
GEOIP,RU,DIRECT
GEOIP,NETFLIX,proxy
```

### Repository structure

```
.
├── go.mod
├── cmd/
│   └── txt2mmdb/
│       └── txt2mmdb.go        # dat → text → mmdb converter
└── .github/
    └── workflows/
        └── build-mmdb.yml     # GitHub Actions pipeline
```

### Download

The latest `geoip-country.mmdb` is always available at a permanent link:

```
https://github.com/kartazon/dat2mmdb/releases/latest/download/geoip-country.mmdb
```

### Run via GitHub Actions

Go to **Actions → build-mmdb → Run workflow**.

The resulting `geoip-country.mmdb` will be published as a GitHub Release asset and available as a workflow artifact.

### Run locally

```bash
# Download geoip.dat
curl -fsSL https://raw.githubusercontent.com/Ground-Zerro/Geo-Aggregator/main/geodat/geoip_GA.dat -o geoip.dat

# Install v2fly geoip tool
go install github.com/v2fly/geoip@latest

# Convert dat → text files
geoip -c config.json

# Convert text → mmdb
mkdir -p output
go run ./cmd/txt2mmdb -indir ./text_output -out output/geoip-country.mmdb
```

### Output MMDB schema

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

### License

MIT

---

## Русский

Конвертирует `geoip.dat` (формат v2fly/v2ray) из репозитория [Ground-Zerro/Geo-Aggregator](https://github.com/Ground-Zerro/Geo-Aggregator) в файл `.mmdb`, совместимый с MaxMind GeoIP. Подходит для использования в [Stash](https://stash.wiki) и других инструментах, поддерживающих **MaxMind GeoIP format**.

### Как это работает

1. `geoip.dat` скачивается автоматически из [Ground-Zerro/Geo-Aggregator](https://raw.githubusercontent.com/Ground-Zerro/Geo-Aggregator/main/geodat/geoip_GA.dat) — загружать файл вручную не нужно.
2. Файл разбирается официальным инструментом [v2fly/geoip](https://github.com/v2fly/geoip) в текстовые файлы с CIDR-диапазонами по тегам.
3. `cmd/txt2mmdb` читает эти файлы и записывает `.mmdb` в формате `GeoLite2-Country` с помощью [`mmdbwriter`](https://github.com/maxmind/mmdbwriter).
4. GitHub Actions запускает этот пайплайн автоматически каждые 12 часов или по ручному триггеру.

В mmdb экспортируются **все теги** из `geoip.dat` — как коды стран (`CN`, `RU`, `US`), так и пользовательские категории (`NETFLIX` и др.). Это позволяет использовать произвольные теги в правилах `GEOIP` в Stash:

```
GEOIP,CN,DIRECT
GEOIP,RU,DIRECT
GEOIP,NETFLIX,proxy
```

### Структура репозитория

```
.
├── go.mod
├── cmd/
│   └── txt2mmdb/
│       └── txt2mmdb.go        # конвертер dat → text → mmdb
└── .github/
    └── workflows/
        └── build-mmdb.yml     # пайплайн GitHub Actions
```

### Скачать

Актуальный `geoip-country.mmdb` всегда доступен по постоянной ссылке:

```
https://github.com/kartazon/dat2mmdb/releases/latest/download/geoip-country.mmdb
```

### Запуск через GitHub Actions

Перейди в **Actions → build-mmdb → Run workflow**.

Готовый `geoip-country.mmdb` будет опубликован как asset GitHub Release и доступен как артефакт workflow.

### Локальный запуск

```bash
# Скачать geoip.dat
curl -fsSL https://raw.githubusercontent.com/Ground-Zerro/Geo-Aggregator/main/geodat/geoip_GA.dat -o geoip.dat

# Установить v2fly geoip
go install github.com/v2fly/geoip@latest

# Конвертировать dat → текстовые файлы
geoip -c config.json

# Конвертировать текст → mmdb
mkdir -p output
go run ./cmd/txt2mmdb -indir ./text_output -out output/geoip-country.mmdb
```

### Схема записи MMDB

Каждая IP-сеть хранится в формате, совместимом с GeoLite2-Country:

```json
{
  "country": {
    "iso_code": "RU",
    "names": { "en": "Russia" }
  },
  "registered_country": {
    "iso_code": "RU",
    "names": { "en": "Russia" }
  }
}
```

### Лицензия

MIT

---

> 🤖 This README was generated with the assistance of [Perplexity AI](https://www.perplexity.ai).
> 
> 🤖 Этот README сформирован с помощью [Perplexity AI](https://www.perplexity.ai).
