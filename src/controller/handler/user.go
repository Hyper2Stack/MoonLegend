package handler

import (
    "encoding/json"
    "net/http"

    "controller/model"
)

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
    if u == nil || u.GetPassword() != hashPassword(in.Password) {
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

    LoginUserVars[r].ResetPassword(hashPassword(in.Password))
}

// PUT /api/v1/user/reset-key
//
func ResetKey(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        Key string `json:"key"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    LoginUserVars[r].ResetKey(in.Key)
}

// GET /api/v1/users/{user_name}
//
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(UserVars[r])
}

// POST /api/v1/user
func AddUser(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    in := struct {
        Name         string `json:"name"`
        DisplayName  string `json:"displayName"`
        Key          string `json:"key"`
        Email        string `json:"email"`
    }{}

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, RequestBodyDecodeError, http.StatusBadRequest)
        return
    }

    if len(in.Name) == 0 {
        http.Error(w, RequestBodyError, http.StatusBadRequest)
        return
    }

    u := new(model.User)
    u.Name = in.Name
    u.DisplayName = in.DisplayName
    u.Key = in.Key
    u.Email = in.Email
    u.Save()

    w.WriteHeader(http.StatusCreated)
}
