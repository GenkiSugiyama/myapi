package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/GenkiSugiyama/myapi/apperrors"
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
		err = apperrors.ReqBodyDecodeFailed.Wrap(err, "bad request body")
		apperrors.ErrorHandler(w, r, err)
		return
	}

	comment, err := c.service.PostCommentService(reqComment)
	if err != nil {
		apperrors.ErrorHandler(w, r, err)
		return
	}

	if err = json.NewEncoder(w).Encode(comment); err != nil {
		err = apperrors.StructEncodeFailed.Wrap(err, "bad struct")
		apperrors.ErrorHandler(w, r, err)
		return
	}
}
