package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

// allowedTags defines which non-ISO tags are written to the mmdb and their priority.
// ISO 2-char country codes are always included with priority 0.
// Higher priority wins on overlapping networks — written last in sorted order.
var allowedTags = map[string]int{
	"RU-WHITELIST": 100,
	"RE-FILTER":    90,
	"TELEGRAM":     80,
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
	"SD": "Sudan", "SE": "Sweden", "SG": "Singapore", "SH": "Saint Helens", "SI": "Slovenia",
	"SJ": "Svalbard and Jan Mayen", "SK": "Slovakia", "SL": "Sierra Leone", "SM": "San Marino", "SN": "Senegal",
	"SO": "Somalia", "SR": "Suriname", "SS": "South Sudan", "ST": "Sao Tome and Principe", "SV": "El Salvador",
	"SX": "Sint Maarten", "SY": "Syria", "SZ": "Eswatini", "TC": "Turks and Caicos Islands", "TD": "Chad",
	"TF": "French Southern Territories", "TG": "Togo", "TH": "Thailand", "TJ": "Tajikistan", "TK": "Tokelau",
	"TL": "Timor-Leste", "TM": "Turkmenistan", "TN": "Tunisia", "TO": "Tonga", "TR": "Turkey",
	"TT": "Trinidad and Tobago", "TV": "Tuvalu", "TW": "Taiwan", "TZ": "Tanzania", "UA": "Ukraine",
	"UG": "Uganda", "UM": "United States Minor Outlying Islands", "US": "United States", "UY": "Uruguay",
	"UZ": "Uzbekistan", "VA": "Vatican City", "VC": "Saint Vincent and the Grenadines", "VE": "Venezuela",
	"VG": "British Virgin Islands", "VI": "U.S. Virgin Islands", "VN": "Vietnam", "VU": "Vanuatu",
	"WF": "Wallis and Futuna", "WS": "Samoa", "YE": "Yemen", "YT": "Mayotte",
	"ZA": "South Africa", "ZM": "Zambia", "ZW": "Zimbabwe",
}

// tagEntry holds a tag name, its file path and resolved priority.
type tagEntry struct {
	tag      string
	path     string
	priority int
}

func main() {
	inDir := flag.String("indir", "./text_output", "directory with per-tag text files produced by v2fly/geoip")
	outFile := flag.String("out", "geoip-country.mmdb", "output mmdb path")
	flag.Parse()

	writer, err := mmdbwriter.New(mmdbwriter.Options{
		DatabaseType:            "GeoLite2-Country",
		IPVersion:               6,
		IncludeReservedNetworks: true,
		Languages:               []string{"en"},
		Description: map[string]string{
			"en": "Country database generated from v2fly geoip.dat",
		},
	})
	if err != nil {
		log.Fatalf("create writer: %v", err)
	}

	entries, err := os.ReadDir(*inDir)
	if err != nil {
		log.Fatalf("read dir %s: %v", *inDir, err)
	}

	// Collect accepted tags and assign priorities.
	var isoTags []tagEntry
	var customTags []tagEntry

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		tag := strings.ToUpper(strings.TrimSuffix(name, filepath.Ext(name)))

		// Skip reserved/private ranges.
		if tag == "PRIVATE" || tag == "LAN" {
			continue
		}

		if len(tag) == 2 {
			// Standard ISO country code — priority 0.
			isoTags = append(isoTags, tagEntry{tag: tag, path: filepath.Join(*inDir, name), priority: 0})
		} else if p, ok := allowedTags[tag]; ok {
			// Allowed custom tag with explicit priority.
			customTags = append(customTags, tagEntry{tag: tag, path: filepath.Join(*inDir, name), priority: p})
		}
		// All other multi-char tags are silently ignored.
	}

	// Sort custom tags ascending by priority so highest-priority is written last
	// and wins on overlapping networks (Insert overwrites previous value).
	sortByPriority(customTags)

	// Write order: ISO countries first (priority 0), then custom tags low→high.
	all := append(isoTags, customTags...)

	totalInserted := 0
	totalSkipped := 0

	for _, e := range all {
		ins, skip := processFile(writer, e.path, e.tag)
		totalInserted += ins
		totalSkipped += skip
		fmt.Fprintf(os.Stderr, "  tag=%-20s inserted=%d skipped=%d\n", e.tag, ins, skip)
	}

	if totalInserted == 0 {
		log.Fatalf("no networks inserted — check that %s contains text files from v2fly/geoip", *inDir)
	}

	out, err := os.Create(*outFile)
	if err != nil {
		log.Fatalf("create output: %v", err)
	}
	defer out.Close()

	if _, err := writer.WriteTo(out); err != nil {
		log.Fatalf("write mmdb: %v", err)
	}

	fmt.Fprintf(os.Stderr, "done: inserted=%d skipped=%d output=%s\n", totalInserted, totalSkipped, *outFile)
}

// sortByPriority sorts tagEntries ascending by priority (lowest first).
func sortByPriority(tags []tagEntry) {
	for i := 1; i < len(tags); i++ {
		for j := i; j > 0 && tags[j].priority < tags[j-1].priority; j-- {
			tags[j], tags[j-1] = tags[j-1], tags[j]
		}
	}
}

func processFile(writer *mmdbwriter.Tree, path, tag string) (inserted, skipped int) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("open %s: %v (skipping)", path, err)
		return
	}
	defer f.Close()

	rec := countryRecord(tag)
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

func countryRecord(iso string) mmdbtype.Map {
	name := countryNames[iso]
	if name == "" {
		name = iso
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
