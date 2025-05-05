package controllers

import (
	"backend/db/socmed"
	"backend/utils"
	"context"
	"encoding/json"
	"net/http"

	"github.com/guregu/null/v5"
)

type CommentController struct {
	queries *socmed.Queries
}

func NewCommentController(queries *socmed.Queries) *CommentController {
	return &CommentController{
		queries: queries,
	}
}

var comment struct {
	CommentUuid    string `json:"comment_uuid"`
	PostUuid       string `json:"post_uuid"`
	UserUuid       string `json:"user_uuid"`
	CommentContent string `json:"comment_content"`
}

func (cc *CommentController) CreateComment(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid create request body", http.StatusBadRequest)
		return
	}

	params := socmed.CreateCommentParams{
		CommentUuid:    utils.GenerateUUID(),
		PostUuid:       null.StringFrom(comment.PostUuid),
		UserUuid:       null.StringFrom(comment.UserUuid),
		CommentContent: null.StringFrom(comment.CommentContent),
	}

	if err := cc.queries.CreateComment(context.Background(), params); err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, map[string]string{"message": "Comment created successfully"})
}

func (cc *CommentController) GetComments(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid get request body", http.StatusBadRequest)
		return
	}

	comments, err := cc.queries.ShowComments(context.Background(), null.StringFrom(comment.PostUuid))
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, comments)
}

func (cc *CommentController) UpdateComment(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid update request body", http.StatusBadRequest)
		return
	}

	params := socmed.UpdateCommentParams{
		CommentUuid:    null.StringFrom(comment.CommentUuid).String,
		CommentContent: null.StringFrom(comment.CommentContent),
	}

	if err := cc.queries.UpdateComment(context.Background(), params); err != nil {
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Comment updated successfully"})
}

func (cc *CommentController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid delete request body", http.StatusBadRequest)
		return
	}

	if err := cc.queries.DeleteComment(context.Background(), comment.PostUuid); err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Comment deleted successfully"})
}
