package domain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Export ...
type Export struct {
	URL string `json:"url"`
}

// AssignedURL will return assigned url to one csv
func (e *Export) AssignedURL(appKey, id string) (string, error) {
	client := &http.Client{}
	url := fmt.Sprintf(
		"%v/%v/url?appKey=%v",
		os.Getenv("EXPORT_URL"),
		id,
		appKey,
	)

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(e); err != nil {
		return "", err
	}

	return e.URL, err
}

// UserAudience ...
type UserAudience struct {
	Email    string
	Birthday string
	Telefone string
}

// ExportDiff ...
type ExportDiff struct {
	IsNew bool
	IsOld bool
	Data  UserAudience
}

type Fn func()
