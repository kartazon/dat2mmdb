package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

var countryNames = map[string]string{
	"AD": "Andorra", "AE": "United Arab Emirates", "AF": "Afghanistan", "AG": "Antigua and Barbuda",
	"AI": "Anguilla", "AL": "Albania", "AM": "Armenia", "AO": "Angola", "AQ": "Antarctica",
	"AR": "Argentina", "AS": "American Samoa", "AT": "Austria", "AU": "Australia", "AW": "Aruba",
	"AX": "Aland Islands", "AZ": "Azerbaijan", "BA": "Bosnia and Herzegovina", "BB": "Barbados",
	"BD": "Bangladesh", "BE": "Belgium", "BF": "Burkina Faso", "BG": "Bulgaria", "BH": "Bahrain",
	"BI": "Burundi", "BJ": "Benin", "BL": "Saint Barthelemy", "BM": "Bermuda", "BN": "Brunei",
	"BO": "Bolivia", "BQ": "Bonaire Sint Eustatius and Saba", "BR": "Brazil", "BS": "Bahamas",
	"BT": "Bhutan", "BV": "Bouvet Island", "BW": "Botswana", "BY": "Belarus", "BZ": "Belize",
	"CA": "Canada", "CC": "Cocos Islands", "CD": "Congo DRC", "CF": "Central African Republic",
	"CG": "Congo", "CH": "Switzerland", "CI": "Cote dIvoire", "CK": "Cook Islands", "CL": "Chile",
	"CM": "Cameroon", "CN": "China", "CO": "Colombia", "CR": "Costa Rica", "CU": "Cuba",
	"CV": "Cabo Verde", "CW": "Curacao", "CX": "Christmas Island", "CY": "Cyprus", "CZ": "Czechia",
	"DE": "Germany", "DJ": "Djibouti", "DK": "Denmark", "DM": "Dominica", "DO": "Dominican Republic",
	"DZ": "Algeria", "EC": "Ecuador", "EE": "Estonia", "EG": "Egypt", "EH": "Western Sahara",
	"ER": "Eritrea", "ES": "Spain", "ET": "Ethiopia", "FI": "Finland", "FJ": "Fiji",
	"FK": "Falkland Islands", "FM": "Micronesia", "FO": "Faroe Islands", "FR": "France", "GA": "Gabon",
	"GB": "United Kingdom", "GD": "Grenada", "GE": "Georgia", "GF": "French Guiana", "GG": "Guernsey",
	"GH": "Ghana", "GI": "Gibraltar", "GL": "Greenland", "GM": "Gambia", "GN": "Guinea",
	"GP": "Guadeloupe", "GQ": "Equatorial Guinea", "GR": "Greece",
	"GS": "South Georgia", "GT": "Guatemala", "GU": "Guam", "GW": "Guinea-Bissau", "GY": "Guyana",
	"HK": "Hong Kong", "HM": "Heard Island", "HN": "Honduras", "HR": "Croatia", "HT": "Haiti",
	"HU": "Hungary", "ID": "Indonesia", "IE": "Ireland", "IL": "Israel", "IM": "Isle of Man",
	"IN": "India", "IO": "British Indian Ocean Territory", "IQ": "Iraq", "IR": "Iran", "IS": "Iceland",
	"IT": "Italy", "JE": "Jersey", "JM": "Jamaica", "JO": "Jordan", "JP": "Japan",
	"KE": "Kenya", "KG": "Kyrgyzstan", "KH": "Cambodia", "KI": "Kiribati", "KM": "Comoros",
	"KN": "Saint Kitts and Nevis", "KP": "North Korea", "KR": "South Korea", "KW": "Kuwait",
	"KY": "Cayman Islands", "KZ": "Kazakhstan", "LA": "Laos", "LB": "Lebanon", "LC": "Saint Lucia",
	"LI": "Liechtenstein", "LK": "Sri Lanka", "LR": "Liberia", "LS": "Lesotho", "LT": "Lithuania",
	"LU": "Luxembourg", "LV": "Latvia", "LY": "Libya", "MA": "Morocco", "MC": "Monaco",
	"MD": "Moldova", "ME": "Montenegro", "MF": "Saint Martin", "MG": "Madagascar",
	"MH": "Marshall Islands", "MK": "North Macedonia", "ML": "Mali", "MM": "Myanmar",
	"MN": "Mongolia", "MO": "Macao", "MP": "Northern Mariana Islands", "MQ": "Martinique",
	"MR": "Mauritania", "MS": "Montserrat", "MT": "Malta", "MU": "Mauritius", "MV": "Maldives",
	"MW": "Malawi", "MX": "Mexico", "MY": "Malaysia", "MZ": "Mozambique", "NA": "Namibia",
	"NC": "New Caledonia", "NE": "Niger", "NF": "Norfolk Island", "NG": "Nigeria", "NI": "Nicaragua",
	"NL": "Netherlands", "NO": "Norway", "NP": "Nepal", "NR": "Nauru", "NU": "Niue",
	"NZ": "New Zealand", "OM": "Oman", "PA": "Panama", "PE": "Peru", "PF": "French Polynesia",
	"PG": "Papua New Guinea", "PH": "Philippines", "PK": "Pakistan", "PL": "Poland",
	"PM": "Saint Pierre and Miquelon", "PN": "Pitcairn", "PR": "Puerto Rico", "PS": "Palestine",
	"PT": "Portugal", "PW": "Palau", "PY": "Paraguay", "QA": "Qatar", "RE": "Reunion",
	"RO": "Romania", "RS": "Serbia", "RU": "Russia", "RW": "Rwanda", "SA": "Saudi Arabia",
	"SB": "Solomon Islands", "SC": "Seychelles", "SD": "Sudan", "SE": "Sweden", "SG": "Singapore",
	"SH": "Saint Helena", "SI": "Slovenia", "SJ": "Svalbard and Jan Mayen", "SK": "Slovakia",
	"SL": "Sierra Leone", "SM": "San Marino", "SN": "Senegal", "SO": "Somalia", "SR": "Suriname",
	"SS": "South Sudan", "ST": "Sao Tome and Principe", "SV": "El Salvador", "SX": "Sint Maarten",
	"SY": "Syria", "SZ": "Eswatini", "TC": "Turks and Caicos Islands", "TD": "Chad",
	"TF": "French Southern Territories", "TG": "Togo", "TH": "Thailand", "TJ": "Tajikistan",
	"TK": "Tokelau", "TL": "Timor-Leste", "TM": "Turkmenistan", "TN": "Tunisia", "TO": "Tonga",
	"TR": "Turkey", "TT": "Trinidad and Tobago", "TV": "Tuvalu", "TW": "Taiwan", "TZ": "Tanzania",
	"UA": "Ukraine", "UG": "Uganda", "UM": "US Minor Outlying Islands", "US": "United States",
	"UY": "Uruguay", "UZ": "Uzbekistan", "VA": "Vatican City",
	"VC": "Saint Vincent and the Grenadines", "VE": "Venezuela", "VG": "British Virgin Islands",
	"VI": "US Virgin Islands", "VN": "Vietnam", "VU": "Vanuatu", "WF": "Wallis and Futuna",
	"WS": "Samoa", "YE": "Yemen", "YT": "Mayotte", "ZA": "South Africa", "ZM": "Zambia",
	"ZW": "Zimbabwe",
}

func main() {
	inDir := flag.String("indir", "./tmp", "directory with per-tag .txt files from v2fly/geoip")
	outFile := flag.String("out", "geoip-country.mmdb", "output mmdb file path")
	cfgFile := flag.String("config", "./mmdb-config.ini", "path to mmdb-config.ini")
	flag.Parse()

	// 1. Discover all available tags from indir
	entries, err := os.ReadDir(*inDir)
	if err != nil {
		log.Fatalf("read dir %s: %v", *inDir, err)
	}

	type fileEntry struct {
		tag  string
		path string
	}
	fileMap := map[string]string{} // tag → file path
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if filepath.Ext(name) != ".txt" {
			continue
		}
		tag := strings.ToUpper(strings.TrimSuffix(name, ".txt"))
		fileMap[tag] = filepath.Join(*inDir, name)
	}

	datTags := make([]string, 0, len(fileMap))
	for tag := range fileMap {
		datTags = append(datTags, tag)
	}
	sort.Strings(datTags)

	// 2. Load and resolve config
	tagConfigs, err := loadConfig(*cfgFile, datTags)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if len(tagConfigs) == 0 {
		log.Fatalf("no tags to export — MMDB would be empty, aborting")
	}

	// 3. Create MMDB writer
	writer, err := mmdbwriter.New(mmdbwriter.Options{
		DatabaseType:            "GeoLite2-Country",
		IPVersion:               6,
		IncludeReservedNetworks: true,
		Languages:               []string{"en"},
		Description: map[string]string{
			"en": "GeoIP database generated from v2fly geoip.dat",
		},
	})
	if err != nil {
		log.Fatalf("create mmdb writer: %v", err)
	}

	// 4. Write tags in priority order (ascending → highest written last, wins on overlap)
	fmt.Fprintln(os.Stderr, "--- Export plan (ascending priority, highest wins) ---")
	for _, tc := range tagConfigs {
		path, ok := fileMap[tc.name]
		if !ok {
			// Already warned in loadConfig for non-ISO tags; ISO tags may be absent
			continue
		}
		ins, skip := processFile(writer, path, tc.name)
		fmt.Fprintf(os.Stderr, "  [priority=%3d] tag=%-25s inserted=%d skipped=%d\n",
			tc.priority, tc.name, ins, skip)
	}

	// 5. Write output
	out, err := os.Create(*outFile)
	if err != nil {
		log.Fatalf("create output file: %v", err)
	}
	defer out.Close()

	if _, err := writer.WriteTo(out); err != nil {
		log.Fatalf("write mmdb: %v", err)
	}

	fmt.Fprintf(os.Stderr, "done → %s\n", *outFile)
}

func processFile(writer *mmdbwriter.Tree, path, tag string) (inserted, skipped int) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("WARN: open %s: %v — skipping tag %s", path, err, tag)
		return
	}
	defer f.Close()

	rec := buildRecord(tag)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		cidr := normalizeCIDR(line)
		if cidr == "" {
			skipped++
			continue
		}
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			skipped++
			continue
		}
		if err := writer.Insert(network, rec); err != nil {
			skipped++
			continue
		}
		inserted++
	}
	return
}

// buildRecord creates the MMDB record for a tag.
// ISO 2-letter tags get a full GeoLite2-Country-compatible structure.
// Special tags (TELEGRAM, RU-WHITELIST, etc.) use iso_code = tag, name = tag.
func buildRecord(tag string) mmdbtype.Map {
	iso := tag
	name := tag
	if n, ok := countryNames[tag]; ok {
		name = n
	}
	names := mmdbtype.Map{"en": mmdbtype.String(name)}
	return mmdbtype.Map{
		"country": mmdbtype.Map{
			"iso_code": mmdbtype.String(iso),
			"names":    names,
		},
		"registered_country": mmdbtype.Map{
			"iso_code": mmdbtype.String(iso),
			"names":    names,
		},
	}
}

func normalizeCIDR(s string) string {
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
