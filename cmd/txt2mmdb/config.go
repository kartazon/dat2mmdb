package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// tagConfig holds the resolved priority for a real tag (after REGION expansion).
type tagConfig struct {
	name     string // real tag name as it appears in dat (uppercase)
	priority int
}

// loadConfig reads mmdb-config.ini and returns the resolved tag configs.
//
// Rules:
//  1. Lines starting with # or empty are ignored.
//  2. REGION is a pseudo-tag — it is expanded to all 2-letter ISO tags found in datTags.
//  3. priority = 0 (or tag absent) → tag is not exported.
//  4. Duplicate priorities → warn + resolve alphabetically (lower name wins).
//  5. Tag in config not found in datTags → warn, skip.
//  6. If config is empty or all priorities are 0 → fallback: REGION=1 implicitly.
func loadConfig(path string, datTags []string) ([]tagConfig, error) {
	raw, err := parseINI(path)
	if err != nil {
		return nil, err
	}

	// Build ISO set from datTags
	isoSet := map[string]bool{}
	datSet := map[string]bool{}
	for _, t := range datTags {
		u := strings.ToUpper(t)
		datSet[u] = true
		if len(u) == 2 {
			isoSet[u] = true
		}
	}

	// Check for empty/all-zero config → implicit REGION=1
	allZero := true
	for _, v := range raw {
		if v > 0 {
			allZero = false
			break
		}
	}
	if allZero {
		fmt.Fprintln(os.Stderr, "WARN: config is empty or all priorities are 0 — exporting countries only (REGION=1 implicit)")
		raw = map[string]int{"REGION": 1}
	}

	// Expand REGION → all ISO tags
	resolved := map[string]int{}
	regionPriority, hasRegion := raw["REGION"]
	if hasRegion && regionPriority > 0 {
		for iso := range isoSet {
			resolved[iso] = regionPriority
		}
	}

	// Process non-REGION tags
	for tag, prio := range raw {
		if tag == "REGION" {
			continue
		}
		if prio == 0 {
			continue
		}
		if !datSet[tag] {
			fmt.Fprintf(os.Stderr, "WARN: tag %q is in config but not found in dat — skipping\n", tag)
			continue
		}
		resolved[tag] = prio
	}

	if len(resolved) == 0 {
		fmt.Fprintln(os.Stderr, "WARN: no tags to export after config processing — output MMDB will be empty")
	}

	// Detect duplicate priorities and warn (collect stats, not per-IP)
	// Build priority → []tag map for collision detection
	prioToTags := map[int][]string{}
	for tag, prio := range resolved {
		prioToTags[prio] = append(prioToTags[prio], tag)
	}
	for prio, tags := range prioToTags {
		if len(tags) > 1 {
			sort.Strings(tags)
			fmt.Fprintf(os.Stderr, "WARN: priority %d is shared by %v — alphabetically first (%s) will win on conflict\n",
				prio, tags, tags[0])
		}
	}

	// Convert to slice and sort: ascending priority (lowest first),
	// ties broken alphabetically (so alphabetically-first is written LAST = wins).
	// Wait — we write ascending, highest written last wins. For ties at same priority,
	// alphabetically FIRST must win → it must be written LAST among ties.
	// So within the same priority, sort DESCENDING by name (last written = first alpha).
	var result []tagConfig
	for tag, prio := range resolved {
		result = append(result, tagConfig{name: tag, priority: prio})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].priority != result[j].priority {
			return result[i].priority < result[j].priority // ascending priority
		}
		// Same priority: sort descending by name so alphabetically-first is written last
		return result[i].name > result[j].name
	})

	return result, nil
}

// parseINI reads a simple KEY = VALUE file, ignoring comments and blank lines.
// Returns map[UPPERCASE_KEY]int_value.
func parseINI(path string) (map[string]int, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		// No config file → treat as empty
		return map[string]int{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open config %s: %w", path, err)
	}
	defer f.Close()

	out := map[string]int{}
	lineNum := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Strip inline comment
		if idx := strings.Index(line, "#"); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "WARN: config line %d invalid (expected KEY = VALUE): %q\n", lineNum, line)
			continue
		}
		key := strings.ToUpper(strings.TrimSpace(parts[0]))
		valStr := strings.TrimSpace(parts[1])
		val, err := strconv.Atoi(valStr)
		if err != nil || val < 0 {
			fmt.Fprintf(os.Stderr, "WARN: config line %d: value %q must be non-negative integer — skipping\n", lineNum, valStr)
			continue
		}
		out[key] = val
	}
	return out, scanner.Err()
}
