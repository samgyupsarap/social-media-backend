package controllers

import (
	"backend/db/socmed"
	"backend/utils"
	"context"
	"encoding/json"
	"net/http"

	"github.com/guregu/null/v5"
)

type PostController struct {
	queries *socmed.Queries
}

func NewPostController(queries *socmed.Queries) *PostController {
	return &PostController{
		queries: queries,
	}
}

var post struct {
	PostUuid    string `json:"post_uuid"`
	UserUuid    string `json:"user_uuid"`
	PostContent string `json:"post_content"`
	PostTags    string `json:"tags"`
	Likes       int    `json:"likes"`
}

func (pc *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid create request body", http.StatusBadRequest)
		return
	}

	params := socmed.CreatePostParams{
		PostUuid:    utils.GenerateUUID(),
		UserUuid:    null.StringFrom(post.UserUuid),
		PostContent: null.StringFrom(post.PostContent),
		PostTags:    null.StringFrom(post.PostTags),
		Likes:       null.Int32From(int32(post.Likes)),
	}

	if err := pc.queries.CreatePost(context.Background(), params); err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, map[string]string{"message": "Post created successfully"})
}

func (pc *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := pc.queries.ShowPosts(context.Background())
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, posts)
}

func (pc *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid update request body", http.StatusBadRequest)
		return
	}

	params := socmed.UpdatePostParams{
		PostContent: null.StringFrom(post.PostContent),
		PostTags:    null.StringFrom(post.PostTags),
	}

	if err := pc.queries.UpdatePost(context.Background(), params); err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Post updated successfully"})

}

func (pc *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := pc.queries.DeletePost(context.Background(), post.PostUuid); err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Post deleted successfully"})
}

func (pc *PostController) UpdateLike(w http.ResponseWriter, r *http.Request) {

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params := socmed.UpdatePostLikesParams{
		Likes:    null.Int32From(int32(post.Likes)),
		PostUuid: post.PostUuid,
	}

	if err := pc.queries.UpdatePostLikes(context.Background(), params); err != nil {
		http.Error(w, "Failed to update like", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Like updated successfully"})
}
