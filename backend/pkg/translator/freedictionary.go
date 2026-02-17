package translator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultFreeDictionaryBaseURL = "https://freedictionaryapi.com"

// WordDetail 单词详细信息
type WordDetail struct {
	Word     string                 `json:"word"`
	Phonetic string                 `json:"phonetic"`
	Meaning  map[string]interface{} `json:"meaning"`
	Example  string                 `json:"example"`
}

// Translator 翻译接口
type Translator interface {
	Translate(word string) (*WordDetail, error)
}

// FreeDictionaryTranslator Free Dictionary API 实现
type FreeDictionaryTranslator struct {
	BaseURL string
	client  *http.Client
}

type freeDictionaryEntry struct {
	Word      string `json:"word"`
	Phonetic  string `json:"phonetic"`
	Phonetics []struct {
		Text string `json:"text"`
	} `json:"phonetics"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string `json:"definition"`
			Example    string `json:"example"`
		} `json:"definitions"`
	} `json:"meanings"`
}

type freeDictionaryError struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// NewFreeDictionaryTranslator 创建 Free Dictionary 翻译器
func NewFreeDictionaryTranslator(baseURL string) *FreeDictionaryTranslator {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = defaultFreeDictionaryBaseURL
	}
	return &FreeDictionaryTranslator{
		BaseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Translate 翻译单词
func (t *FreeDictionaryTranslator) Translate(word string) (*WordDetail, error) {
	normalized := strings.TrimSpace(word)
	if normalized == "" {
		return nil, fmt.Errorf("word is empty")
	}

	entries, err := t.requestEntries(normalized)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("word not found: %s", normalized)
	}

	entry := entries[0]
	phonetic := entry.Phonetic
	if phonetic == "" {
		for _, p := range entry.Phonetics {
			if strings.TrimSpace(p.Text) != "" {
				phonetic = p.Text
				break
			}
		}
	}

	definitions := make([]map[string]string, 0)
	var example string
	for _, meaning := range entry.Meanings {
		for _, def := range meaning.Definitions {
			text := strings.TrimSpace(def.Definition)
			if text == "" {
				continue
			}
			item := map[string]string{"text": text}
			if strings.TrimSpace(meaning.PartOfSpeech) != "" {
				item["pos"] = meaning.PartOfSpeech
			}
			definitions = append(definitions, item)
			if example == "" && strings.TrimSpace(def.Example) != "" {
				example = def.Example
			}
		}
	}
	if len(definitions) == 0 {
		return nil, fmt.Errorf("no definitions found for word: %s", normalized)
	}

	return &WordDetail{
		Word:     normalized,
		Phonetic: phonetic,
		Meaning: map[string]interface{}{
			"definitions": definitions,
		},
		Example: example,
	}, nil
}

func (t *FreeDictionaryTranslator) requestEntries(word string) ([]freeDictionaryEntry, error) {
	escaped := url.PathEscape(word)
	endpoints := []string{
		t.BaseURL + "/api/v1/entries/en/" + escaped,
		t.BaseURL + "/api/v2/entries/en/" + escaped,
	}

	var lastErr error
	for _, endpoint := range endpoints {
		entries, handled, err := t.tryEndpoint(endpoint)
		if handled {
			if err != nil {
				lastErr = err
				continue
			}
			return entries, nil
		}
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("failed to fetch word from free dictionary api")
}

func (t *FreeDictionaryTranslator) tryEndpoint(endpoint string) ([]freeDictionaryEntry, bool, error) {
	resp, err := t.client.Get(endpoint)
	if err != nil {
		return nil, true, fmt.Errorf("failed to call free dictionary api: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, true, fmt.Errorf("failed to read free dictionary response: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		var apiErr freeDictionaryError
		if json.Unmarshal(body, &apiErr) == nil && strings.TrimSpace(apiErr.Title) != "" {
			return nil, true, fmt.Errorf("word not found")
		}
		return nil, true, fmt.Errorf("word not found")
	}

	if resp.StatusCode == http.StatusNotImplemented || resp.StatusCode == http.StatusMethodNotAllowed {
		return nil, false, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, true, fmt.Errorf("free dictionary api status: %d", resp.StatusCode)
	}

	var entries []freeDictionaryEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, true, fmt.Errorf("failed to parse free dictionary response: %w", err)
	}
	return entries, true, nil
}
