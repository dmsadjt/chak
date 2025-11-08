package search

import(
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type BraveManager struct {
	szApiUrl string
	szApiKey string
}

func NewBraveManager(szApiKey string) *BraveManager {
	return &BraveManager{
		szApiUrl: "https://api.search.brave.com/res/v1/web/search",
		szApiKey: szApiKey,
	}
}

type braveResponse struct {
	Web struct {
		Results []struct {
			SzTitle string `json:"title"`
			SzSnippet string `json:"description"`
			SzURL string `json:"url"`
		} `json:"results"`
	}`json:"web"`
} 
 
func (braveMgr *BraveManager) Search(SzQuery string) ([]SearchResultData, error) {
	szSearchUrl := fmt.Sprintf("%s?q=%s&count=5", braveMgr.szApiUrl, url.QueryEscape(SzQuery))

	req, err := http.NewRequest("GET", szSearchUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Chak-Server/1.0")
	req.Header.Set("X-Subscription-Token", braveMgr.szApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw braveResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var searchResults []SearchResultData
	for _, topic := range raw.Web.Results {
		searchResults = append(searchResults, SearchResultData {
			SzTitle: topic.SzTitle,
			SzSnippet: topic.SzSnippet,
			SzURL: topic.SzURL,
		})
	}

	return searchResults, nil
}
