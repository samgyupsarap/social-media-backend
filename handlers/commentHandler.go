package handlers

import (
	"backend/controllers"
	"backend/db/socmed"
	"net/http"
)

type CommentHandlerController struct {
	commentController *controllers.CommentController
}

func NewCommentHandlerController(queries *socmed.Queries) *CommentHandlerController {
	return &CommentHandlerController{
		commentController: controllers.NewCommentController(queries),
	}
}

func (h *CommentHandlerController) CommentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.commentController.GetComments(w, r)
	case http.MethodPost:
		h.commentController.CreateComment(w, r)
	case http.MethodPatch:
		h.commentController.UpdateComment(w, r)
	case http.MethodDelete:
		h.commentController.DeleteComment(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
