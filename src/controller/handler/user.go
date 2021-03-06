package handler

import (
    "encoding/json"
    "net/http"

    "controller/model"
)

// POST /api/v1/signup
//
func Signup(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        Username string `json:"username"`
        Password string `json:"password"`
        Email    string `json:"email"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    // TBD, should verify username & password

    if model.GetUserByName(in.Username) != nil {
        w.WriteHeader(http.StatusConflict)
        return
    }

    u := new(model.User)
    u.Name = in.Username
    u.Email = in.Email
    u.Save()

    u.ResetPassword(in.Password)
}

// POST /api/v1/login
//
func Login(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    u := model.GetUserByName(in.Username)
    if u == nil || !u.VerifyPassword(in.Password) {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    out := struct {
        Key string `json:"key"`
    }{}

    out.Key = encodeUserToken(in.Username)
    json.NewEncoder(w).Encode(out)
}

// GET /api/v1/user
//
func GetMyProfile(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(LoginUserVars[r])
}

// PUT /api/v1/user/reset-password
//
func ResetPassword(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        Password string `json:"password"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    // TBD, should verify password

    LoginUserVars[r].ResetPassword(in.Password)
}

// PUT /api/v1/user/reset-key
//
func ResetKey(w http.ResponseWriter, r *http.Request) {
    LoginUserVars[r].ResetKey()
}
