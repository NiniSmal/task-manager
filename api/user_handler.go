package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gitlab.com/nina8884807/task-manager/entity"
)

type UserHandler struct {
	service UserService
	appHost string
}

func NewUserHandler(u UserService, appHost string) *UserHandler {
	return &UserHandler{
		service: u,
		appHost: appHost,
	}
}

type UserService interface {
	CreateUser(ctx context.Context, login, password, photo string) error
	Login(ctx context.Context, login, password string) (uuid.UUID, error)
	Verification(ctx context.Context, verificationCode string, verification bool) error
	ResendVerificationCode(ctx context.Context, email string) error
	SendAnAbsenceLetter(ctx context.Context, intervalTime string) error
	UploadPhoto(ctx context.Context, imageURL string) error
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	err = u.service.CreateUser(ctx, user.Email, user.Password, user.Photo)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	sessionID, err := u.service.Login(ctx, user.Email, user.Password)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	l := r.Context().Value("logger").(*slog.Logger)

	l.Info(fmt.Sprintf("login OK session_id: %s", sessionID))

	cookie := http.Cookie{
		Name:  "session_id",
		Value: sessionID.String(),
		Path:  "/",
		// Domain:     u.appHost,
		Expires:    time.Now().Add(time.Hour * 24 * 30),
		RawExpires: "",
		MaxAge:     86400 * 7,
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteNoneMode,
		Raw:        "",
		Unparsed:   nil,
	}
	http.SetCookie(w, &cookie)
}
func (u *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:  "session_id",
		Value: "-",
		Path:  "/",
		// Domain:     u.appHost,
		Expires:    time.Now().Add(-time.Hour),
		RawExpires: "",
		MaxAge:     -1,
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteNoneMode,
		Raw:        "",
		Unparsed:   nil,
	}
	http.SetCookie(w, &cookie)

}

func (u *UserHandler) Verification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user entity.User

	code := r.URL.Query().Get("code")
	if code == "" {
		HandlerError(ctx, w, errors.New("code is empty"))
		return
	}

	user.VerificationCode = code
	user.Verification = true

	err := u.service.Verification(ctx, user.VerificationCode, user.Verification)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	http.Redirect(w, r, "https://tm.anaxita.ru/login", http.StatusSeeOther)
}

func (u *UserHandler) RepeatRequestVerification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
	err = u.service.ResendVerificationCode(ctx, user.Email)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
}

type imageJSON struct {
	Base64 string `json:"base64"`
}

func (u *UserHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	file, _, err := r.FormFile("file")
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	var base64Encoding string

	mimeType := http.DetectContentType(content)
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}
	base64Encoding += base64.StdEncoding.EncodeToString(content)

	err = u.service.UploadPhoto(ctx, base64Encoding)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
}
