package main

import "strings"

// vowelTones lists the four tone-marked forms of each pinyin vowel, from
// tone 1 through tone 4. Tone 5 is written on the unmarked vowel.
var vowelTones = map[rune][]rune{
	'a': {'ā', 'á', 'ǎ', 'à'},
	'e': {'ē', 'é', 'ě', 'è'},
	'i': {'ī', 'í', 'ǐ', 'ì'},
	'o': {'ō', 'ó', 'ǒ', 'ò'},
	'u': {'ū', 'ú', 'ǔ', 'ù'},
	'ü': {'ǖ', 'ǘ', 'ǚ', 'ǜ'},
}

// markedPinyin renders a stored numbered reading for the Telegram message.
// Syllables that are already tone-marked, or that are not numbered pinyin,
// are copied through unchanged.
func markedPinyin(reading string) string {
	reading = strings.TrimSpace(reading)
	if reading == "" {
		return ""
	}

	var message strings.Builder
	var syllable []rune
	flush := func() {
		if len(syllable) == 0 {
			return
		}
		message.WriteString(markSyllable(string(syllable)))
		syllable = syllable[:0]
	}
	for _, r := range reading {
		if isPinyinRune(r) {
			syllable = append(syllable, r)
			continue
		}
		flush()
		message.WriteRune(r)
	}
	flush()
	return message.String()
}

func markSyllable(token string) string {
	for _, r := range token {
		if isToneMark(r) {
			return token
		}
	}

	normalized := strings.ToLower(token)
	normalized = strings.ReplaceAll(normalized, "u:", "ü")
	normalized = strings.ReplaceAll(normalized, "v", "ü")
	runes := []rune(normalized)
	if len(runes) < 2 {
		return token
	}

	toneDigit := runes[len(runes)-1]
	if toneDigit < '1' || toneDigit > '5' {
		return token
	}
	base := runes[:len(runes)-1]
	for _, r := range base {
		if !isPinyinLetter(r) {
			return token
		}
	}

	tone := int(toneDigit - '0')
	if tone == 5 {
		return string(base)
	}

	index := toneVowelIndex(base)
	if index < 0 {
		return token
	}
	marked := vowelTones[base[index]][tone-1]
	base[index] = marked
	return string(base)
}

func toneVowelIndex(base []rune) int {
	for _, vowel := range []rune{'a', 'e', 'o'} {
		for i, r := range base {
			if r == vowel {
				return i
			}
		}
	}
	for i := len(base) - 1; i >= 0; i-- {
		switch base[i] {
		case 'i', 'u', 'ü':
			return i
		}
	}
	return -1
}

func isPinyinLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || r == 'ü'
}

func isPinyinRune(r rune) bool {
	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == ':' || r == 'ü' || r == 'Ü' {
		return true
	}
	return isToneMark(r)
}

func isToneMark(r rune) bool {
	switch r {
	case 'ā', 'á', 'ǎ', 'à', 'ē', 'é', 'ě', 'è',
		'ī', 'í', 'ǐ', 'ì', 'ō', 'ó', 'ǒ', 'ò',
		'ū', 'ú', 'ǔ', 'ù', 'ǖ', 'ǘ', 'ǚ', 'ǜ',
		'Ā', 'Á', 'Ǎ', 'À', 'Ē', 'É', 'Ě', 'È',
		'Ī', 'Í', 'Ǐ', 'Ì', 'Ō', 'Ó', 'Ǒ', 'Ò',
		'Ū', 'Ú', 'Ǔ', 'Ù', 'Ǖ', 'Ǘ', 'Ǚ', 'Ǜ':
		return true
	default:
		return false
	}
}
