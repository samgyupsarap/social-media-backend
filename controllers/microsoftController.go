package controllers

import (
	"backend/db/socmed"
	"backend/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/guregu/null/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type MicrosoftController struct {
	config  *oauth2.Config
	queries *socmed.Queries
}

type MicrosoftUser struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"userPrincipalName"`
}

func NewMicrosoftController(queries *socmed.Queries) *MicrosoftController {
	config := &oauth2.Config{
		ClientID:     os.Getenv("MICROSOFT_CLIENT_ID"),
		ClientSecret: os.Getenv("MICROSOFT_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("MICROSOFT_REDIRECT_URL"),
		Scopes: []string{
			"User.Read",
			"offline_access",
		},
		Endpoint: microsoft.AzureADEndpoint(os.Getenv("MICROSOFT_TENANT_ID")),
	}
	return &MicrosoftController{
		config:  config,
		queries: queries,
	}
}

func (mc *MicrosoftController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	url := mc.config.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (mc *MicrosoftController) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	msToken, err := mc.config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := mc.config.Client(r.Context(), msToken)
	resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	var user MicrosoftUser
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		return
	}

	// Here you would typically:
	// 1. Check if user exists in your database
	// 2. Create user if they don't exist
	// 3. Generate a session token
	// 4. Return the token to the client

	existingUser, err := mc.queries.SelectUserByEmail(r.Context(), null.StringFrom(user.Email))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not registered", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	appUser := socmed.User{
		UserUuid: existingUser.UserUuid,
		Email:    existingUser.Email,
		UserName: existingUser.UserName,
		FullName: existingUser.FullName,
	}

	token, err := utils.GenerateToken(appUser)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"message":  "Login successful",
		"msToken":  msToken,
		"token":    token,
		"uuid":     appUser.UserUuid,
		"email":    appUser.Email,
		"username": appUser.UserName,
		"fullName": appUser.FullName,
	})

	frontendURL := os.Getenv("FRONTEND_URL")
	redirectURL := fmt.Sprintf("%s/auth/callback?"+
		"token=%s&"+
		"uuid=%s&"+
		"email=%s&"+
		"username=%s&"+
		"fullName=%s&"+
		"message=%s",
		frontendURL,
		url.QueryEscape(token),
		url.QueryEscape(appUser.UserUuid),
		url.QueryEscape(appUser.Email.String),
		url.QueryEscape(appUser.UserName.String),
		url.QueryEscape(appUser.FullName.String),
		url.QueryEscape("Login successful"),
	)

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
