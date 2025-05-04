package controllers

import (
	"backend/db/socmed"
	"backend/utils"
	"context"
	"net/http"

	"io"

	"github.com/guregu/null/v5"
)

type UserController struct {
	queries *socmed.Queries
}

func NewUserController(queries *socmed.Queries) *UserController {
	return &UserController{
		queries: queries,
	}
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	fullName := r.FormValue("full_name")
	email := r.FormValue("email")
	userName := r.FormValue("user_name")
	password := r.FormValue("password")

	file, _, err := r.FormFile("profile_picture")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error retrieving profile picture", http.StatusBadRequest)
		return
	}

	var profilePictures [][]byte
	if file != nil {
		defer file.Close()

		profilePicture, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading profile picture", http.StatusInternalServerError)
			return
		}
		profilePictures = append(profilePictures, profilePicture)
	}

	params := socmed.CreateUserParams{
		UserUuid:       utils.GenerateUUID(),
		FullName:       null.StringFrom(fullName),
		Email:          null.StringFrom(email),
		UserName:       null.StringFrom(userName),
		Password:       null.StringFrom(password),
		ProfilePicture: profilePictures,
	}

	if err := uc.queries.CreateUser(context.Background(), params); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
	})
}
func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	userUUID := queryParams.Get("uuid")

	if userUUID == "" {
		http.Error(w, "UUID is required", http.StatusBadRequest)
		return
	}

	user, err := uc.queries.SelectUserByUUID(context.Background(), userUUID)
	if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, user)
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	userUUID := r.FormValue("uuid")
	fullName := r.FormValue("full_name")
	email := r.FormValue("email")
	userName := r.FormValue("user_name")
	password := r.FormValue("password")

	file, _, err := r.FormFile("profile_picture")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error retrieving profile picture", http.StatusBadRequest)
		return
	}

	var profilePictures [][]byte
	if file != nil {
		defer file.Close()

		profilePicture, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading profile picture", http.StatusInternalServerError)
			return
		}
		profilePictures = append(profilePictures, profilePicture)
	}

	params := socmed.UpdateUserParams{
		UserUuid:       userUUID,
		FullName:       null.StringFrom(fullName),
		Email:          null.StringFrom(email),
		UserName:       null.StringFrom(userName),
		Password:       null.StringFrom(password),
		ProfilePicture: profilePictures,
	}

	if err := uc.queries.UpdateUser(context.Background(), params); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "User updated successfully",
	})
}
