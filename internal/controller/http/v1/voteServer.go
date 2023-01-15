package http_v1

import (
	"beta/internal/adapters/db/postgresql"
	"beta/internal/domain/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	Db postgresql.VoteStorage
}

// PostVote adds a vote from JSON received in the request body.
func (s *Server) PostVote(c *gin.Context) {
	var newVote models.RequestVote

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&newVote); err != nil {
		return
	}

	var responseOutside = models.Response{Result: "ok"}
	c.IndentedJSON(http.StatusOK, responseOutside)
}
