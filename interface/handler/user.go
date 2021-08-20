package handler

import (
	"net/http"

	"go-chi-ddd/resource/request"
	"go-chi-ddd/usecase"
)

type User struct {
	userUseCase usecase.IUser
}

func NewUser(uuc usecase.IUser) User {
	return User{
		userUseCase: uuc,
	}
}

func (u User) Create(w http.ResponseWriter, r *http.Request) error {
	var req request.UserCreate

	if !bind(w, r, &req) {
		return nil
	}

	id, err := u.userUseCase.Create(newCtx(), &req)
	if err != nil {
		return err
	}

	err = writeJson(w, http.StatusCreated, id)
	if err != nil {
		return err
	}
	return nil
}

func (u User) ResetPasswordRequest(w http.ResponseWriter, r *http.Request) error {
	var req request.UserResetPasswordRequest

	if !bind(w, r, &req) {
		return nil
	}

	res, err := u.userUseCase.ResetPasswordRequest(newCtx(), &req)
	if err != nil {
		return err
	}

	err = writeJson(w, http.StatusOK, res)
	if err != nil {
		return err
	}
	return nil
}

func (u User) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	var req request.UserResetPassword

	if !bind(w, r, &req) {
		return nil
	}

	err := u.userUseCase.ResetPassword(newCtx(), &req)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (u User) Login(w http.ResponseWriter, r *http.Request) error {
	var req request.UserLogin

	if !bind(w, r, &req) {
		return nil
	}

	res, err := u.userUseCase.Login(newCtx(), &req)
	if err != nil {
		return err
	}

	if res == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	err = writeJson(w, http.StatusOK, res)
	if err != nil {
		return err
	}
	return nil
}

func (u User) RefreshToken(w http.ResponseWriter, r *http.Request) error {
	res, err := u.userUseCase.RefreshToken(r.URL.Query().Get("refresh_token"))
	if err != nil {
		return err
	}

	if res == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	err = writeJson(w, http.StatusCreated, res)
	if err != nil {
		return err
	}
	return nil
}
