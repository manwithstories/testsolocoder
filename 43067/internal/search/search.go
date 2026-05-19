package search

import (
	"strings"

	"snippetbox/internal/logger"
	"snippetbox/internal/models"
)

func containsFuzzy(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)

	if len(substr) == 0 {
		return true
	}

	if strings.Contains(s, substr) {
		return true
	}

	j := 0
	for i := 0; i < len(s) && j < len(substr); i++ {
		if s[i] == substr[j] {
			j++
		}
	}
	return j == len(substr)
}

func containsAllTags(snippetTags, queryTags []string) bool {
	tagSet := make(map[string]bool)
	for _, t := range snippetTags {
		tagSet[strings.ToLower(t)] = true
	}

	for _, qt := range queryTags {
		if !tagSet[strings.ToLower(qt)] {
			return false
		}
	}
	return true
}

func SearchSnippets(snippets []models.Snippet, query models.SearchQuery) []models.Snippet {
	var results []models.Snippet

	searchFields := query.Fields
	if len(searchFields) == 0 {
		searchFields = []string{"title", "tags", "language", "content", "description"}
	}

	for _, snippet := range snippets {
		matched := true

		if query.Language != "" {
			if !strings.EqualFold(snippet.Language, query.Language) {
				matched = false
			}
		}

		if matched && len(query.Tags) > 0 {
			if !containsAllTags(snippet.Tags, query.Tags) {
				matched = false
			}
		}

		if matched && query.Keyword != "" {
			keywordMatched := false
			for _, field := range searchFields {
				switch strings.ToLower(field) {
				case "title":
					if containsFuzzy(snippet.Title, query.Keyword) {
						keywordMatched = true
					}
				case "content":
					if containsFuzzy(snippet.Content, query.Keyword) {
						keywordMatched = true
					}
				case "description":
					if containsFuzzy(snippet.Description, query.Keyword) {
						keywordMatched = true
					}
				case "language":
					if containsFuzzy(snippet.Language, query.Keyword) {
						keywordMatched = true
					}
				case "tags":
					for _, tag := range snippet.Tags {
						if containsFuzzy(tag, query.Keyword) {
							keywordMatched = true
							break
						}
					}
				}
				if keywordMatched {
					break
				}
			}
			if !keywordMatched {
				matched = false
			}
		}

		if matched {
			results = append(results, snippet)
		}
	}

	logger.Info("Search completed: %d results for query '%s'", len(results), query.Keyword)
	return results
}
