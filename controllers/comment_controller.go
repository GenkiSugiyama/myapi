package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/GenkiSugiyama/myapi/controllers/services"
	"github.com/GenkiSugiyama/myapi/models"
)

type CommentController struct {
	service services.CommentServicer
}

func NewCommentController(s services.CommentServicer) *CommentController {
	return &CommentController{service: s}
}

func (c *CommentController) PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	var reqComment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&reqComment); err != nil {
		http.Error(w, "fail to decode json\n", http.StatusBadRequest)
		return
	}

	comment, err := c.service.PostCommentService(reqComment)
	if err != nil {
		http.Error(w, "failed internal exec\n", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(comment); err != nil {
		http.Error(w, "failed to encode json\n", http.StatusInternalServerError)
		return
	}
}
