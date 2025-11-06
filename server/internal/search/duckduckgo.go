package search

import(
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type DuckDuckGoManager struct {
	szApiUrl string
}

func NewDuckDuckGoManager() *DuckDuckGoManager {
	return &DuckDuckGoManager{
		szApiUrl: "https://api.duckduckgo.com/",
	}
}

type duckduckgoResponse struct {
	RelatedTopics []struct {
		SzText string `json:"Text"`
		SzFirstURL string
	} `json:"RelatedTopics"`
}

func (duckMgr *DuckDuckGoManager) Search(SzQuery string) ([]SearchResultData, error) {
	szSearchUrl := fmt.Sprintf("%s?q=%s&format=json", duckMgr.szApiUrl, url.QueryEscape(SzQuery))

	req, err := http.NewRequest("GET", szSearchUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Chak-Server/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw duckduckgoResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	results := make([]SearchResultData, 0, len(raw.RelatedTopics))
	for _, topic := range raw.RelatedTopics {
		results = append(results, SearchResultData {
			SzTitle: topic.SzText,
			SzSnippet: topic.SzText,
			SzURL: topic.SzFirstURL,
		})
	}

	if err != nil {
		return nil, err
	}

	return results, nil
}
