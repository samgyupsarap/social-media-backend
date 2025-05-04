package models

type UserInput struct {
    FullName string `json:"full_name" form:"full_name"`
    Email    string `json:"email" form:"email"`
    UserName string `json:"user_name" form:"user_name"`
    Password string `json:"password" form:"password"`
}