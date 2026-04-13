package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "net"
    "os"
    "sort"
    "strings"

    "github.com/maxmind/mmdbwriter"
    "github.com/maxmind/mmdbwriter/mmdbtype"
)

type record struct {
    Country mmdbtype.Map `maxminddb:"country"`
}

var countryNames = map[string]string{
    "AD": "Andorra", "AE": "United Arab Emirates", "AF": "Afghanistan", "AG": "Antigua and Barbuda",
    "AI": "Anguilla", "AL": "Albania", "AM": "Armenia", "AO": "Angola", "AQ": "Antarctica",
    "AR": "Argentina", "AS": "American Samoa", "AT": "Austria", "AU": "Australia", "AW": "Aruba",
    "AX": "Aland Islands", "AZ": "Azerbaijan", "BA": "Bosnia and Herzegovina", "BB": "Barbados",
    "BD": "Bangladesh", "BE": "Belgium", "BF": "Burkina Faso", "BG": "Bulgaria", "BH": "Bahrain",
    "BI": "Burundi", "BJ": "Benin", "BL": "Saint Barthelemy", "BM": "Bermuda", "BN": "Brunei",
    "BO": "Bolivia", "BQ": "Bonaire, Sint Eustatius and Saba", "BR": "Brazil", "BS": "Bahamas",
    "BT": "Bhutan", "BV": "Bouvet Island", "BW": "Botswana", "BY": "Belarus", "BZ": "Belize",
    "CA": "Canada", "CC": "Cocos (Keeling) Islands", "CD": "Congo (DRC)", "CF": "Central African Republic",
    "CG": "Congo", "CH": "Switzerland", "CI": "Cote d'Ivoire", "CK": "Cook Islands", "CL": "Chile",
    "CM": "Cameroon", "CN": "China", "CO": "Colombia", "CR": "Costa Rica", "CU": "Cuba",
    "CV": "Cabo Verde", "CW": "Curacao", "CX": "Christmas Island", "CY": "Cyprus", "CZ": "Czechia",
    "DE": "Germany", "DJ": "Djibouti", "DK": "Denmark", "DM": "Dominica", "DO": "Dominican Republic",
    "DZ": "Algeria", "EC": "Ecuador", "EE": "Estonia", "EG": "Egypt", "EH": "Western Sahara",
    "ER": "Eritrea", "ES": "Spain", "ET": "Ethiopia", "FI": "Finland", "FJ": "Fiji",
    "FK": "Falkland Islands", "FM": "Micronesia", "FO": "Faroe Islands", "FR": "France", "GA": "Gabon",
    "GB": "United Kingdom", "GD": "Grenada", "GE": "Georgia", "GF": "French Guiana", "GG": "Guernsey",
    "GH": "Ghana", "GI": "Gibraltar", "GL": "Greenland", "GM": "Gambia", "GN": "Guinea",
    "GP": "Guadeloupe", "GQ": "Equatorial Guinea", "GR": "Greece", "GS": "South Georgia and the South Sandwich Islands",
    "GT": "Guatemala", "GU": "Guam", "GW": "Guinea-Bissau", "GY": "Guyana", "HK": "Hong Kong",
    "HM": "Heard Island and McDonald Islands", "HN": "Honduras", "HR": "Croatia", "HT": "Haiti", "HU": "Hungary",
    "ID": "Indonesia", "IE": "Ireland", "IL": "Israel", "IM": "Isle of Man", "IN": "India",
    "IO": "British Indian Ocean Territory", "IQ": "Iraq", "IR": "Iran", "IS": "Iceland", "IT": "Italy",
    "JE": "Jersey", "JM": "Jamaica", "JO": "Jordan", "JP": "Japan", "KE": "Kenya",
    "KG": "Kyrgyzstan", "KH": "Cambodia", "KI": "Kiribati", "KM": "Comoros", "KN": "Saint Kitts and Nevis",
    "KP": "North Korea", "KR": "South Korea", "KW": "Kuwait", "KY": "Cayman Islands", "KZ": "Kazakhstan",
    "LA": "Laos", "LB": "Lebanon", "LC": "Saint Lucia", "LI": "Liechtenstein", "LK": "Sri Lanka",
    "LR": "Liberia", "LS": "Lesotho", "LT": "Lithuania", "LU": "Luxembourg", "LV": "Latvia",
    "LY": "Libya", "MA": "Morocco", "MC": "Monaco", "MD": "Moldova", "ME": "Montenegro",
    "MF": "Saint Martin", "MG": "Madagascar", "MH": "Marshall Islands", "MK": "North Macedonia", "ML": "Mali",
    "MM": "Myanmar", "MN": "Mongolia", "MO": "Macao", "MP": "Northern Mariana Islands", "MQ": "Martinique",
    "MR": "Mauritania", "MS": "Montserrat", "MT": "Malta", "MU": "Mauritius", "MV": "Maldives",
    "MW": "Malawi", "MX": "Mexico", "MY": "Malaysia", "MZ": "Mozambique", "NA": "Namibia",
    "NC": "New Caledonia", "NE": "Niger", "NF": "Norfolk Island", "NG": "Nigeria", "NI": "Nicaragua",
    "NL": "Netherlands", "NO": "Norway", "NP": "Nepal", "NR": "Nauru", "NU": "Niue",
    "NZ": "New Zealand", "OM": "Oman", "PA": "Panama", "PE": "Peru", "PF": "French Polynesia",
    "PG": "Papua New Guinea", "PH": "Philippines", "PK": "Pakistan", "PL": "Poland", "PM": "Saint Pierre and Miquelon",
    "PN": "Pitcairn", "PR": "Puerto Rico", "PS": "Palestine", "PT": "Portugal", "PW": "Palau",
    "PY": "Paraguay", "QA": "Qatar", "RE": "Reunion", "RO": "Romania", "RS": "Serbia",
    "RU": "Russia", "RW": "Rwanda", "SA": "Saudi Arabia", "SB": "Solomon Islands", "SC": "Seychelles",
    "SD": "Sudan", "SE": "Sweden", "SG": "Singapore", "SH": "Saint Helena", "SI": "Slovenia",
    "SJ": "Svalbard and Jan Mayen", "SK": "Slovakia", "SL": "Sierra Leone", "SM": "San Marino", "SN": "Senegal",
    "SO": "Somalia", "SR": "Suriname", "SS": "South Sudan", "ST": "Sao Tome and Principe", "SV": "El Salvador",
    "SX": "Sint Maarten", "SY": "Syria", "SZ": "Eswatini", "TC": "Turks and Caicos Islands", "TD": "Chad",
    "TF": "French Southern Territories", "TG": "Togo", "TH": "Thailand", "TJ": "Tajikistan", "TK": "Tokelau",
    "TL": "Timor-Leste", "TM": "Turkmenistan", "TN": "Tunisia", "TO": "Tonga", "TR": "Turkey",
    "TT": "Trinidad and Tobago", "TV": "Tuvalu", "TW": "Taiwan", "TZ": "Tanzania", "UA": "Ukraine",
    "UG": "Uganda", "UM": "United States Minor Outlying Islands", "US": "United States", "UY": "Uruguay",
    "UZ": "Uzbekistan", "VA": "Vatican City", "VC": "Saint Vincent and the Grenadines", "VE": "Venezuela",
    "VG": "British Virgin Islands", "VI": "U.S. Virgin Islands", "VN": "Vietnam", "VU": "Vanuatu", "WF": "Wallis and Futuna",
    "WS": "Samoa", "YE": "Yemen", "YT": "Mayotte", "ZA": "South Africa", "ZM": "Zambia", "ZW": "Zimbabwe",
}

func main() {
    inFile := flag.String("in", "geoip.txt", "input text file produced by v2fly/geoip")
    outFile := flag.String("out", "geoip-country.mmdb", "output mmdb path")
    flag.Parse()

    f, err := os.Open(*inFile)
    if err != nil {
        log.Fatalf("open input: %v", err)
    }
    defer f.Close()

    writer, err := mmdbwriter.New(mmdbwriter.Options{
        DatabaseType: "GeoLite2-Country",
        Description: map[string]string{
            "en": "Country database generated from v2fly geoip.dat",
        },
        DisableIPv4Aliasing: false,
        IncludeReservedNetworks: true,
        IPVersion: 6,
        Languages: []string{"en"},
    })
    if err != nil {
        log.Fatalf("create writer: %v", err)
    }

    scanner := bufio.NewScanner(f)
    buf := make([]byte, 0, 1024*1024)
    scanner.Buffer(buf, 10*1024*1024)

    currentTag := ""
    inserted := 0
    skipped := 0
    seen := map[string]struct{}{}

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
            currentTag = strings.ToUpper(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"))
            continue
        }

        if currentTag == "" {
            skipped++
            continue
        }

        if currentTag == "PRIVATE" || currentTag == "LAN" {
            continue
        }

        prefix := normalizeCIDR(line)
        if prefix == "" {
            skipped++
            continue
        }

        key := currentTag + "|" + prefix
        if _, ok := seen[key]; ok {
            continue
        }
        seen[key] = struct{}{}

        _, network, err := net.ParseCIDR(prefix)
        if err != nil {
            skipped++
            continue
        }

        rec := countryRecord(currentTag)
        if err := writer.Insert(network, rec); err != nil {
            skipped++
            continue
        }
        inserted++
    }
    if err := scanner.Err(); err != nil {
        log.Fatalf("scan input: %v", err)
    }

    if inserted == 0 {
        log.Fatalf("no networks inserted, check geoip.txt format")
    }

    out, err := os.Create(*outFile)
    if err != nil {
        log.Fatalf("create output: %v", err)
    }
    defer out.Close()

    if _, err := writer.WriteTo(out); err != nil {
        log.Fatalf("write mmdb: %v", err)
    }

    fmt.Fprintf(os.Stderr, "done: inserted=%d skipped=%d output=%s\n", inserted, skipped, *outFile)
}

func normalizeCIDR(s string) string {
    s = strings.TrimSpace(s)
    if s == "" {
        return ""
    }
    if strings.Contains(s, "/") {
        return s
    }
    ip := net.ParseIP(s)
    if ip == nil {
        return ""
    }
    if ip.To4() != nil {
        return s + "/32"
    }
    return s + "/128"
}

func countryRecord(tag string) mmdbtype.Map {
    iso := strings.ToUpper(tag)
    name := countryNames[iso]
    if name == "" {
        name = iso
    }

    names := mmdbtype.Map{
        "en": mmdbtype.String(name),
    }

    rec := mmdbtype.Map{
        "country": mmdbtype.Map{
            "iso_code": mmdbtype.String(iso),
            "names":    names,
        },
        "registered_country": mmdbtype.Map{
            "iso_code": mmdbtype.String(iso),
            "names":    names,
        },
    }

    return rec
}

func init() {
    keys := make([]string, 0, len(countryNames))
    for k := range countryNames {
        keys = append(keys, k)
    }
    sort.Strings(keys)
}
