package handlers

import (
	"backend/controllers"
	"backend/db/socmed"
	"net/http"
)

type PostHandlerController struct {
	postController *controllers.PostController
}

func NewPostHandlerController(queries *socmed.Queries) *PostHandlerController {
	return &PostHandlerController{
		postController: controllers.NewPostController(queries),
	}
}

func (h *PostHandlerController) PostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.postController.GetPosts(w, r)
	case http.MethodPost:
		h.postController.CreatePost(w, r)
	case http.MethodPatch:
		h.postController.UpdatePost(w, r)
	case http.MethodDelete:
		h.postController.DeletePost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
