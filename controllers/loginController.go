package controllers

import (
	"backend/db/socmed"
	"backend/utils"
	"context"
	"encoding/json"
	"net/http"

	"github.com/guregu/null/v5"
)

type LoginController struct {
	queries *socmed.Queries
}

func NewLoginController(queries *socmed.Queries) *LoginController {
	return &LoginController{
		queries: queries,
	}
}

type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func (lc *LoginController) Login(w http.ResponseWriter, r *http.Request) {

	var input LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := lc.queries.SelectUserByUserName(context.Background(), null.StringFrom(input.UserName))
	if err != nil {
		http.Error(w, "User not found "+err.Error(), http.StatusNotFound)
		return
	}

	if user.Password.String != input.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	mappedUser := socmed.User{
		UserUuid: user.UserUuid,
		UserName: user.UserName,
		Email:    user.Email,
		FullName: user.FullName,
		Password: user.Password,
	}

	token, err := utils.GenerateToken(mappedUser)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"message":   "Login successful",
		"token":     token,
		"user_uuid": user.UserUuid,
		"user_name": user.UserName.String,
		"email":     user.Email.String,
		"full_name": user.FullName.String,
	})
}
