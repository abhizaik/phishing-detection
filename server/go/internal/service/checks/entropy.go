package checks

import (
	"fmt"
	"math"
	"strings"
	"unicode"
)

// DomainRandomnessResult holds metrics and final score
type DomainRandomnessResult struct {
	Domain            string
	Label             string
	Length            int
	Entropy           float64
	EntropyPerChar    float64
	NormalizedEntropy float64

	VowelRatio          float64
	DigitRatio          float64
	UniqueCharRatio     float64
	LongestConsonantRun int
	BigramEnglishiness  float64

	RandomnessScore float64
	IsSuspicious    bool
	Reasons         []string
}

// ---------- Constants & small resources ----------
const maxAlphabetBits = 5.954196310386875 // log2(62)

var commonBigrams = map[string]struct{}{
	"th": {}, "he": {}, "in": {}, "er": {}, "an": {}, "re": {}, "on": {}, "at": {}, "en": {}, "nd": {},
	"ti": {}, "es": {}, "or": {}, "te": {}, "of": {}, "ed": {}, "is": {}, "it": {}, "al": {}, "ar": {},
	"st": {}, "to": {}, "nt": {}, "ng": {}, "se": {}, "ha": {}, "as": {}, "ou": {}, "io": {}, "le": {},
}

// ---------- Helpers ----------

// isAllowedRune checks if rune is a-z, 0-9, or hyphen (we lowercase runes first at call site)
func isAllowedRune(r rune) bool {
	if r == '-' {
		return true
	}
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// shannonEntropy computes Shannon entropy in bits over runes
func shannonEntropy(s string) float64 {
	if s == "" {
		return 0
	}
	freq := make(map[rune]int)
	total := 0
	for _, r := range s {
		freq[r]++
		total++
	}
	var H float64
	for _, count := range freq {
		p := float64(count) / float64(total)
		H -= p * math.Log2(p)
	}
	return H
}

func entropyPerChar(s string) float64 {
	if s == "" {
		return 0
	}
	return shannonEntropy(s) / float64(len([]rune(s)))
}

func normalizedEntropy(s string) float64 {
	epc := entropyPerChar(s)
	return epc / maxAlphabetBits
}

// extractSLD returns the second-level label if available (example.com -> "example")
// keeps behavior simple (does not handle public suffix list). That's intentional — you can swap to PSL if needed.
func extractSLD(domain string) string {
	domain = strings.TrimSpace(strings.ToLower(domain))
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	if i := strings.IndexAny(domain, "/?#"); i >= 0 {
		domain = domain[:i]
	}
	if i := strings.Index(domain, ":"); i >= 0 {
		domain = domain[:i]
	}
	parts := strings.Split(domain, ".")
	if len(parts) >= 2 {
		return parts[len(parts)-2]
	}
	return strings.Join(parts, "")
}

// bigramEnglishiness: fraction of letter-only bigrams that occur in a small English bigram set
func bigramEnglishiness(s string) float64 {
	filtered := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsLetter(r) {
			filtered = append(filtered, unicode.ToLower(r))
		}
	}
	if len(filtered) < 2 {
		return 0.0
	}
	found := 0
	total := 0
	for i := 0; i < len(filtered)-1; i++ {
		bg := string([]rune{filtered[i], filtered[i+1]})
		total++
		if _, ok := commonBigrams[bg]; ok {
			found++
		}
	}
	if total == 0 {
		return 0.0
	}
	return float64(found) / float64(total)
}

func uniqueCharRatio(s string) float64 {
	seen := make(map[rune]struct{})
	total := 0
	for _, r := range s {
		if isAllowedRune(r) {
			seen[r] = struct{}{}
			total++
		}
	}
	if total == 0 {
		return 0
	}
	return float64(len(seen)) / float64(total)
}

func vowelRatio(s string) float64 {
	vowels := "aeiou"
	count := 0
	total := 0
	for _, r := range s {
		if !isAllowedRune(r) {
			continue
		}
		total++
		if strings.ContainsRune(vowels, unicode.ToLower(r)) {
			count++
		}
	}
	if total == 0 {
		return 0
	}
	return float64(count) / float64(total)
}

func digitRatio(s string) float64 {
	count := 0
	total := 0
	for _, r := range s {
		if !isAllowedRune(r) {
			continue
		}
		total++
		if unicode.IsDigit(r) {
			count++
		}
	}
	if total == 0 {
		return 0
	}
	return float64(count) / float64(total)
}

func longestConsonantRunLength(s string) int {
	vowels := "aeiou"
	maxRun := 0
	cur := 0
	for _, r := range s {
		lr := unicode.ToLower(r)
		if !isAllowedRune(lr) {
			cur = 0
			continue
		}
		// digits count as consonant/digit run
		if unicode.IsDigit(lr) || (unicode.IsLetter(lr) && !strings.ContainsRune(vowels, lr)) {
			cur++
		} else {
			if cur > maxRun {
				maxRun = cur
			}
			cur = 0
		}
	}
	if cur > maxRun {
		maxRun = cur
	}
	return maxRun
}

// ---------- Analyzer ----------

// AnalyzeDomainRandomness computes multiple heuristics and a combined randomness score.
// The thresholds and weights are hard-coded here (tweak as needed), but kept grouped and obvious.
func AnalyzeDomainRandomness(domain string) DomainRandomnessResult {
	label := extractSLD(domain)

	// sanitize label: lowercase, keep only a-z0-9-
	var b strings.Builder
	for _, r := range label {
		lr := unicode.ToLower(r)
		if isAllowedRune(lr) {
			b.WriteRune(lr)
		}
	}
	s := b.String()

	// prepare baseline metrics (handle empty/special)
	ent := shannonEntropy(s)
	epc := entropyPerChar(s)
	normE := normalizedEntropy(s)
	vRatio := vowelRatio(s)
	dRatio := digitRatio(s)
	uRatio := uniqueCharRatio(s)
	longRun := longestConsonantRunLength(s)
	bigram := bigramEnglishiness(s)

	// normalize some features
	longRunNorm := math.Min(float64(longRun)/6.0, 1.0) // saturate at 6+

	// weights (grouped for readability)
	weights := struct {
		NormalizedEntropy float64
		VowelInverse      float64
		DigitRatio        float64
		BigramInverse     float64
		UniqueInverse     float64
		LongRun           float64
	}{
		NormalizedEntropy: 0.25,
		VowelInverse:      0.15,
		DigitRatio:        0.15,
		BigramInverse:     0.20,
		UniqueInverse:     0.10,
		LongRun:           0.15,
	}

	// compute weighted randomness score (higher => more random-looking)
	score := 0.0
	score += weights.NormalizedEntropy * normE
	score += weights.VowelInverse * (1.0 - vRatio)
	score += weights.DigitRatio * dRatio
	score += weights.BigramInverse * (1.0 - bigram)
	score += weights.UniqueInverse * (1.0 - uRatio)
	score += weights.LongRun * longRunNorm

	// clamp to [0,1]
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	// Decide suspiciousness and collect reasons (deduplicate later)
	isSusp := false
	reasons := make([]string, 0, 4)
	baseThreshold := 0.50

	// length-adjusted threshold (smooth boost for short labels)
	adjThreshold := baseThreshold
	runeLen := len([]rune(s))
	if runeLen > 0 && runeLen < 6 {
		boost := 0.10 * (float64(6-runeLen) / 5.0) // len=5 -> +0.02, len=1 -> +0.10
		adjThreshold = baseThreshold + boost
	}

	if score > adjThreshold {
		isSusp = true
		reasons = append(reasons, fmt.Sprintf("high randomness score=%.3f>threshold=%.3f", score, adjThreshold))
	}

	// hard/auxiliary rules
	if longRun >= 6 {
		isSusp = true
		reasons = append(reasons, fmt.Sprintf("long consonant/digit run=%d", longRun))
	}
	if dRatio > 0.5 && runeLen >= 6 {
		isSusp = true
		reasons = append(reasons, fmt.Sprintf("high digit ratio=%.2f", dRatio))
	}
	if uRatio < 0.25 && runeLen >= 6 {
		reasons = append(reasons, fmt.Sprintf("low unique-char-ratio=%.2f", uRatio))
	}

	// deduplicate reasons
	uniqueReasons := make([]string, 0, len(reasons))
	seen := make(map[string]struct{}, len(reasons))
	for _, r := range reasons {
		if _, ok := seen[r]; ok {
			continue
		}
		seen[r] = struct{}{}
		uniqueReasons = append(uniqueReasons, r)
	}

	return DomainRandomnessResult{
		Domain:              domain,
		Label:               label,
		Length:              runeLen,
		Entropy:             ent,
		EntropyPerChar:      epc,
		NormalizedEntropy:   normE,
		VowelRatio:          vRatio,
		DigitRatio:          dRatio,
		UniqueCharRatio:     uRatio,
		LongestConsonantRun: longRun,
		BigramEnglishiness:  bigram,
		RandomnessScore:     score,
		IsSuspicious:        isSusp,
		Reasons:             uniqueReasons,
	}
}
