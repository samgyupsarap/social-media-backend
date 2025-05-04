package handlers

import (
	"backend/controllers"
	"backend/db/socmed"
	"net/http"
)

type UserHandlerController struct {
	userController *controllers.UserController
}

func NewUserHandlerController(queries *socmed.Queries) *UserHandlerController {
	return &UserHandlerController{
		userController: controllers.NewUserController(queries),
	}
}

func (h *UserHandlerController) UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.userController.GetUser(w, r)
	case http.MethodPost:
		h.userController.CreateUser(w, r)
	case http.MethodPatch:
		h.userController.UpdateUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
