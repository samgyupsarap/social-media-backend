package controllers

import (
	"backend/db/socmed"
	"backend/utils"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/guregu/null/v5"
	"golang.org/x/oauth2"
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

func generateCodeVerifier() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func generateCodeChallenge(verifier string) string {
	h := sha256.New()
	h.Write([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func NewMicrosoftController(queries *socmed.Queries) *MicrosoftController {
    config := &oauth2.Config{
        ClientID:     os.Getenv("MICROSOFT_CLIENT_ID"),
        ClientSecret: os.Getenv("MICROSOFT_CLIENT_SECRET"),
        RedirectURL:  os.Getenv("MICROSOFT_REDIRECT_URI"),
        Scopes: []string{
            "https://graph.microsoft.com/User.Read",
            "openid",
            "profile",
            "email",
        },
        Endpoint: oauth2.Endpoint{
            AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
            TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
        },
    }
    return &MicrosoftController{
        config:  config,
        queries: queries,
    }
}

func (mc *MicrosoftController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate PKCE verifier and challenge
	verifier := generateCodeVerifier()
	challenge := generateCodeChallenge(verifier)

	// Store verifier in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "code_verifier",
		Value:    verifier,
		Path:     "/",
		MaxAge:   int(time.Hour.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// Add PKCE parameters to auth URL
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	}

	url := mc.config.AuthCodeURL("state", opts...)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (mc *MicrosoftController) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Get stored verifier
	verifierCookie, err := r.Cookie("code_verifier")
	if err != nil {
		http.Error(w, "Verifier not found", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Exchange code for token with PKCE
	msToken, err := mc.config.Exchange(r.Context(), code,
		oauth2.SetAuthURLParam("code_verifier", verifierCookie.Value))
	if err != nil {
		fmt.Printf("Token exchange error: %v\n", err)
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Get user info from Microsoft Graph
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

	// Check if user exists
	existingUser, err := mc.queries.SelectUserByEmail(r.Context(), null.StringFrom(user.Email))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not registered", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Create app user
	appUser := socmed.User{
		UserUuid: existingUser.UserUuid,
		Email:    existingUser.Email,
		UserName: existingUser.UserName,
		FullName: existingUser.FullName,
	}

	// Generate JWT token
	token, err := utils.GenerateToken(appUser)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Redirect to frontend with user data
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
