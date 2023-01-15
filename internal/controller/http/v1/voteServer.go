package http_v1

import (
	"beta/internal/domain/models"
	vote_service "beta/internal/domain/service"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PostVote adds a vote from JSON received in the request body.
func PostVote(c *gin.Context) {
	var newVote models.RequestVote

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&newVote); err != nil {
		return
	}

	go vote_service.GService.AddVote(context.Background(), newVote)

	var responseOutside = models.Response{Result: "ok"}

	c.IndentedJSON(http.StatusOK, responseOutside)
}
