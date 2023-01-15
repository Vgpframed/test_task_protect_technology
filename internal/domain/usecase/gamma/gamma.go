package gamma

import (
	"beta/internal/domain/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendVotesGamma(votes []models.SendVote, url string) {

	fullUrl := fmt.Sprintf("%svoting-stats", url)
	reqBody, err := json.Marshal(&votes)
	if err != nil {
		print(err)
	}
	resp, err := http.Post(fullUrl,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()

}
