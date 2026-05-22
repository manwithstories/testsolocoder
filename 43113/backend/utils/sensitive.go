package utils

import (
	"qa-platform/models"
	"strings"
	"sync"
)

type SensitiveWordFilter struct {
	trie      *TrieNode
	words     []models.SensitiveWord
	mu        sync.RWMutex
}

type TrieNode struct {
	children  map[rune]*TrieNode
	isEnd     bool
	word      string
	replaceTo string
}

var SensitiveFilter *SensitiveWordFilter

func InitSensitiveFilter(words []models.SensitiveWord) {
	SensitiveFilter = &SensitiveWordFilter{
		trie:  &TrieNode{children: make(map[rune]*TrieNode)},
		words: words,
	}
	for _, word := range words {
		SensitiveFilter.AddWord(word.Word, word.ReplaceTo)
	}
}

func (f *SensitiveWordFilter) AddWord(word, replaceTo string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	node := f.trie
	for _, ch := range word {
		if node.children[ch] == nil {
			node.children[ch] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[ch]
	}
	node.isEnd = true
	node.word = word
	node.replaceTo = replaceTo
}

func (f *SensitiveWordFilter) Filter(content string) (string, []string) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	result := []rune(content)
	foundWords := make([]string, 0)

	i := 0
	for i < len(result) {
		node := f.trie
		matchLen := 0
		foundMatch := false
		var matchedWord string
		var replaceTo string

		for j := i; j < len(result); j++ {
			ch := result[j]
			if next, ok := node.children[ch]; ok {
				node = next
				if node.isEnd {
					matchLen = j - i + 1
					foundMatch = true
					matchedWord = node.word
					replaceTo = node.replaceTo
				}
			} else {
				break
			}
		}

		if foundMatch && matchLen > 0 {
			foundWords = append(foundWords, matchedWord)
			replacement := replaceTo
			if replacement == "" {
				replacement = strings.Repeat("*", matchLen)
			}
			result = []rune(string(result[:i]) + replacement + string(result[i+matchLen:]))
			i += len([]rune(replacement))
		} else {
			i++
		}
	}

	return string(result), foundWords
}

func (f *SensitiveWordFilter) Check(content string) (bool, []string) {
	_, foundWords := f.Filter(content)
	return len(foundWords) > 0, foundWords
}
