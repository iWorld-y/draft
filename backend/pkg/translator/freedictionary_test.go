package translator

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFreeDictionaryTranslatorTranslateSuccess(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/entries/en/behavior" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"word":"behavior",
			"entries":[
				{
					"partOfSpeech":"noun",
					"pronunciations":[
						{"type":"respelling","text":"bih-HAYV-yer"},
						{"type":"ipa","text":"/bɪˈheɪvjər/"}
					],
					"senses":[
						{
							"definition":"The way a living creature behaves.",
							"examples":["Her behavior changed over time."]
						}
					]
				}
			]
		}`))
	}))
	defer server.Close()

	tr := NewFreeDictionaryTranslator(server.URL)
	got, err := tr.Translate("behavior")
	if err != nil {
		t.Fatalf("Translate returned error: %v", err)
	}

	if got.Word != "behavior" {
		t.Fatalf("unexpected word: %s", got.Word)
	}
	if got.Phonetic != "/bɪˈheɪvjər/" {
		t.Fatalf("unexpected phonetic: %s", got.Phonetic)
	}
	if got.Example != "Her behavior changed over time." {
		t.Fatalf("unexpected example: %s", got.Example)
	}

	definitions, ok := got.Meaning["definitions"].([]map[string]string)
	if !ok {
		t.Fatalf("definitions has unexpected type: %T", got.Meaning["definitions"])
	}
	if len(definitions) != 1 {
		t.Fatalf("unexpected definitions count: %d", len(definitions))
	}
	if definitions[0]["text"] != "The way a living creature behaves." {
		t.Fatalf("unexpected definition text: %s", definitions[0]["text"])
	}
	if definitions[0]["pos"] != "noun" {
		t.Fatalf("unexpected partOfSpeech: %s", definitions[0]["pos"])
	}
}

func TestFreeDictionaryTranslatorTranslateNotFound(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"title":"No Definitions Found","message":"not found"}`))
	}))
	defer server.Close()

	tr := NewFreeDictionaryTranslator(server.URL)
	_, err := tr.Translate("behavior")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "word not found") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFreeDictionaryTranslatorTranslateNoDefinitions(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"word":"behavior",
			"entries":[
				{
					"partOfSpeech":"noun",
					"pronunciations":[{"type":"ipa","text":"/bɪˈheɪvjər/"}],
					"senses":[{"definition":"   ","examples":["example"]}]
				}
			]
		}`))
	}))
	defer server.Close()

	tr := NewFreeDictionaryTranslator(server.URL)
	_, err := tr.Translate("behavior")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "no definitions found") {
		t.Fatalf("unexpected error: %v", err)
	}
}
