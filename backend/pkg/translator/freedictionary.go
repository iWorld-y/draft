package translator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
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

type freeDictionaryResponse struct {
	Word    string                `json:"word"`
	Entries []freeDictionaryEntry `json:"entries"`
}

type freeDictionaryEntry struct {
	PartOfSpeech   string `json:"partOfSpeech"`
	Pronunciations []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"pronunciations"`
	Senses []struct {
		Definition string   `json:"definition"`
		Examples   []string `json:"examples"`
	} `json:"senses"`
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

	resp, err := t.requestEntries(normalized)
	if err != nil {
		return nil, err
	}
	if len(resp.Entries) == 0 {
		return nil, fmt.Errorf("word not found: %s", normalized)
	}

	phonetic := pickPhonetic(resp.Entries)

	definitions := make([]map[string]string, 0)
	var example string
	for _, entry := range resp.Entries {
		for _, sense := range entry.Senses {
			text := strings.TrimSpace(sense.Definition)
			if text == "" {
				continue
			}
			item := map[string]string{"text": text}
			if strings.TrimSpace(entry.PartOfSpeech) != "" {
				item["pos"] = entry.PartOfSpeech
			}
			definitions = append(definitions, item)
			if example == "" {
				for _, ex := range sense.Examples {
					if strings.TrimSpace(ex) != "" {
						example = ex
						break
					}
				}
			}
		}
	}
	if len(definitions) == 0 {
		return nil, fmt.Errorf("no definitions found for word: %s", normalized)
	}

	wordInResp := strings.TrimSpace(resp.Word)
	if wordInResp == "" {
		wordInResp = normalized
	}

	return &WordDetail{
		Word:     wordInResp,
		Phonetic: phonetic,
		Meaning: map[string]interface{}{
			"definitions": definitions,
		},
		Example: example,
	}, nil
}

func (t *FreeDictionaryTranslator) requestEntries(word string) (*freeDictionaryResponse, error) {
	escaped := url.PathEscape(word)
	endpoint := t.BaseURL + "/api/v1/entries/en/" + escaped
	return t.tryEndpoint(endpoint)
}

func (t *FreeDictionaryTranslator) tryEndpoint(endpoint string) (*freeDictionaryResponse, error) {
	resp, err := t.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to call free dictionary api: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read free dictionary response: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Warnf("free dictionary api returned 404 endpoint=%s body_preview=%q", endpoint, truncateBodyPreview(body))
		var apiErr freeDictionaryError
		if json.Unmarshal(body, &apiErr) == nil && strings.TrimSpace(apiErr.Title) != "" {
			return nil, fmt.Errorf("word not found")
		}
		return nil, fmt.Errorf("word not found")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("free dictionary api status: %d", resp.StatusCode)
	}

	var parsed freeDictionaryResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse free dictionary response: %w", err)
	}
	return &parsed, nil
}

func pickPhonetic(entries []freeDictionaryEntry) string {
	for _, entry := range entries {
		for _, p := range entry.Pronunciations {
			if strings.EqualFold(strings.TrimSpace(p.Type), "ipa") && strings.TrimSpace(p.Text) != "" {
				return p.Text
			}
		}
	}
	for _, entry := range entries {
		for _, p := range entry.Pronunciations {
			if strings.TrimSpace(p.Text) != "" {
				return p.Text
			}
		}
	}
	return ""
}

func truncateBodyPreview(body []byte) string {
	preview := strings.TrimSpace(string(body))
	const maxLen = 256
	if len(preview) <= maxLen {
		return preview
	}
	return preview[:maxLen]
}
