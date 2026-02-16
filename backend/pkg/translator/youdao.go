// pkg/translator/youdao.go
package translator

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

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

// YoudaoTranslator 有道翻译实现
type YoudaoTranslator struct {
	AppKey    string
	AppSecret string
	client    *http.Client
}

// YoudaoResponse 有道 API 响应结构
type YoudaoResponse struct {
	ErrorCode string `json:"errorCode"`
	Query     string `json:"query"`
	Basic     struct {
		Phonetic string   `json:"phonetic"`
		Explains []string `json:"explains"`
	} `json:"basic"`
	Web []struct {
		Key   string   `json:"key"`
		Value []string `json:"value"`
	} `json:"web"`
}

// NewYoudaoTranslator 创建有道翻译器
func NewYoudaoTranslator(appKey, appSecret string) *YoudaoTranslator {
	return &YoudaoTranslator{
		AppKey:    appKey,
		AppSecret: appSecret,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Translate 翻译单词
func (t *YoudaoTranslator) Translate(word string) (*WordDetail, error) {
	if t.AppKey == "" || t.AppSecret == "" {
		return nil, fmt.Errorf("youdao app key or secret is empty")
	}

	// 1. 构建请求参数
	salt := strconv.FormatInt(time.Now().Unix(), 10)
	sign := t.generateSign(word, salt)

	params := url.Values{}
	params.Add("q", word)
	params.Add("from", "en")
	params.Add("to", "zh-CHS")
	params.Add("appKey", t.AppKey)
	params.Add("salt", salt)
	params.Add("sign", sign)

	// 2. 发送请求
	apiURL := "https://openapi.youdao.com/api?" + params.Encode()
	resp, err := t.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call youdao api: %w", err)
	}
	defer resp.Body.Close()

	// 3. 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result YoudaoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if result.ErrorCode != "0" {
		return nil, fmt.Errorf("youdao api error: %s", result.ErrorCode)
	}

	// 4. 转换为内部数据结构
	return &WordDetail{
		Word:     word,
		Phonetic: result.Basic.Phonetic,
		Meaning:  t.parseMeaning(result.Basic.Explains),
	}, nil
}

// generateSign 生成签名
func (t *YoudaoTranslator) generateSign(word, salt string) string {
	str := t.AppKey + word + salt + t.AppSecret
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}

// parseMeaning 解析释义
func (t *YoudaoTranslator) parseMeaning(explains []string) map[string]interface{} {
	definitions := make([]map[string]string, 0, len(explains))
	for _, exp := range explains {
		definitions = append(definitions, map[string]string{
			"text": exp,
		})
	}
	return map[string]interface{}{
		"definitions": definitions,
	}
}
