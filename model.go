package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

type Sense struct {
	Gloss     string `json:"gloss"`
	SortOrder int    `json:"sort_order"`
}

type VocabularyEntry struct {
	ID          string  `json:"id"`
	Simplified  string  `json:"simplified"`
	Traditional string  `json:"traditional"`
	Reading     string  `json:"reading"`
	Senses      []Sense `json:"senses"`
}

func (entry VocabularyEntry) Meaning() string {
	if len(entry.Senses) == 0 {
		return ""
	}

	lowest := entry.Senses[0]
	for _, sense := range entry.Senses[1:] {
		if sense.SortOrder < lowest.SortOrder {
			lowest = sense
		}
	}
	return lowest.Gloss
}

type selector struct {
	random *rand.Rand
}

func newSelector(random *rand.Rand) *selector {
	if random == nil {
		panic("nil random source")
	}
	return &selector{random: random}
}

func (s *selector) Select(entries []VocabularyEntry, limit int) []VocabularyEntry {
	if limit <= 0 || len(entries) == 0 {
		return nil
	}

	shuffled := append([]VocabularyEntry(nil), entries...)
	s.random.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	if len(shuffled) > limit {
		shuffled = shuffled[:limit]
	}
	return shuffled
}

func formatDailyMessage(entries []VocabularyEntry) string {
	if len(entries) == 0 {
		return "No confirmed words to review today."
	}

	var message strings.Builder
	message.WriteString("Today's words\n\n")
	for index, entry := range entries {
		fmt.Fprintf(&message, "%d. %s", index+1, entry.Simplified)
		traditional := strings.TrimSpace(entry.Traditional)
		if traditional != "" && traditional != entry.Simplified {
			fmt.Fprintf(&message, "（%s）", traditional)
		}
		message.WriteByte('\n')
		message.WriteString(entry.Reading)
		message.WriteByte('\n')
		message.WriteString(entry.Meaning())
		if index < len(entries)-1 {
			message.WriteString("\n\n")
		}
	}
	return message.String()
}
