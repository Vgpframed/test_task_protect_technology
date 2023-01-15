package gamma

import (
	"beta/internal/domain/models"
	"bytes"
	"encoding/json"
	"fmt"
	lg "gitlab.satel.eyevox.ru/satel_vks/jaeger_tracer/log"
	"go.uber.org/zap"
	"net/http"
)

func SendVotesGamma(votes []models.SendVote, url string, logger lg.Factory) {

	fullUrl := fmt.Sprintf("%svoting-stats", url)
	reqBody, err := json.Marshal(&votes)
	if err != nil {
		logger.Bg().Error("sending request to SendVotesGamma", zap.Error(err))
	}
	resp, err := http.Post(fullUrl,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		logger.Bg().Error("sending request to SendVotesGamma", zap.Error(err))
	}
	defer resp.Body.Close()

}
