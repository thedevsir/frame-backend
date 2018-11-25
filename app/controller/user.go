package controller

import (
	"io"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsir/frame-backend/app/repository"
	"github.com/thedevsir/frame-backend/config"
	"github.com/thedevsir/frame-backend/services/auth"
	"github.com/thedevsir/frame-backend/services/errors"
	j "github.com/thedevsir/frame-backend/services/jwt"
	"github.com/thedevsir/frame-backend/services/mail"
	"github.com/thedevsir/frame-backend/services/request"
	r "github.com/thedevsir/frame-backend/services/response"
	"github.com/thedevsir/frame-backend/services/storage"
	"github.com/thedevsir/frame-backend/services/validation"
)

var (
	_SendVerficationMail = mail.SendVerficationMail
	_SendResetMail       = mail.SendResetMail
)

type (
	SignupShcema struct {
		Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
		Password string `json:"password" validate:"required,min=8,max=50"`
		Email    string `json:"email" validate:"required,email"`
	}
	ResendShcema struct {
		Email string `json:"email" validate:"required,email"`
	}
	VerificationShcema struct {
		Token string `json:"token" validate:"required"`
	}
	SigninShcema struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	ForgotShcema struct {
		Email string `json:"email" validate:"required,email"`
	}
	ResetShcema struct {
		Token    string `json:"token" validate:"required"`
		Password string `json:"password" validate:"required,min=8,max=50"`
	}
	ChangePasswordShcema struct {
		Password string `json:"password" validate:"required,min=8,max=50"`
	}
	ChangeUsernameShcema struct {
		Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	}
	ChangeEmailShcema struct {
		Email string `json:"email" validate:"required,email"`
	}
)

// Signup godoc
// @Summary Create an account
// @Tags user
// @Accept json
// @Produce json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Param email body string true "Email"
// @Success 201 {object} response.Message
// @Router /users/signup [post]
func Signup(c echo.Context) (err error) {

	params := new(SignupShcema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	if _, err = repository.CheckUsername(params.Username); err != nil {
		return err
	}

	if _, err = repository.CheckEmail(params.Email); err != nil {
		return err
	}

	user, err := repository.CreateUser(params.Username, params.Password, params.Email)
	if err != nil {
		return err
	}

	token, err := mail.MakeEmailToken("verify", user.GetId().Hex(), params.Username, params.Email, []byte(config.SigningKey))
	if err == nil {
		go _SendVerficationMail(params.Username, params.Email, token)
	}

	return errors.ErrCreated
}

// Resend godoc
// @Summary Resend email verfication
// @Tags user
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Success 200 {object} response.Message
// @Router /users/signup/resend [post]
func Resend(c echo.Context) (err error) {

	params := new(ResendShcema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	if user, err := repository.CheckEmail(params.Email); err != nil {

		if user.IsEmailVerified {
			return errors.ErrAccountVerified
		}

		token, err := mail.MakeEmailToken("verify", user.GetId().Hex(), user.Username, params.Email, []byte(config.SigningKey))
		if err == nil {
			go _SendVerficationMail(user.Username, params.Email, token)
		}

		return errors.ErrSuccess
	}

	return errors.ErrUserEmailNotFound
}

// Verification godoc
// @Summary Activation
// @Tags user
// @Accept json
// @Produce json
// @Param token body string true "Token"
// @Success 200 {object} response.Message
// @Router /users/signup/verification [post]
func Verification(c echo.Context) (err error) {

	params := new(VerificationShcema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	data, err := j.ParseEmailToken(params.Token, []byte(config.SigningKey))
	if err != nil {
		return err
	}

	if data.Action != "verify" {
		return errors.ErrAccessDenied
	}

	if err = repository.UserActivation(data.UserID); err != nil {
		return err
	}

	return errors.ErrSuccess
}

// Signin godoc
// @Summary User signin
// @Tags user
// @Accept json
// @Produce json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 201 {object} response.Message
// @Router /users/signin [post]
func Signin(c echo.Context) (err error) {

	ip := c.RealIP()
	userAgent := c.Request().Header.Get("User-Agent")

	params := new(SigninShcema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	if err = repository.CheckAbuse(ip, params.Username); err != nil {
		return err
	}

	user, err := repository.FindUserByCredentials(params.Username, params.Password)
	if err != nil {
		repository.SubmitAttempt(ip, params.Username)
		return err
	}

	userID := user.Id.Hex()
	SID, uuid, err := repository.SessionCreate(ip, userID, userAgent)
	if err != nil {
		return err
	}

	token := &auth.UserToken{
		uuid,
		SID,
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	jwtToken, err := token.Create([]byte(config.SigningKey))
	if err != nil {
		return err
	}

	return r.CustomErrorWithHeader(http.StatusOK, "Success", map[string]string{echo.HeaderAuthorization: jwtToken}, c)
}

// Forgot godoc
// @Summary Forgot password
// @Tags user
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Success 200 {object} response.Message
// @Router /users/signin/forgot [post]
func Forgot(c echo.Context) (err error) {

	params := new(ForgotShcema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	if user, err := repository.CheckEmail(params.Email); err != nil {

		token, err := mail.MakeEmailToken("reset", user.GetId().Hex(), user.Username, params.Email, []byte(config.SigningKey))
		if err == nil {
			go _SendResetMail(user.Username, params.Email, token)
		}

		return errors.ErrSuccess
	}

	return errors.ErrUserEmailNotFound
}

// Reset godoc
// @Summary Reset password
// @Tags user
// @Accept json
// @Produce json
// @Param token body string true "Token"
// @Param password body string true "Password"
// @Success 200 {object} response.Message
// @Router /users/signin/reset [put]
func Reset(c echo.Context) (err error) {

	params := new(ResetShcema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	data, err := j.ParseEmailToken(params.Token, []byte(config.SigningKey))
	if err != nil {
		return err
	}

	if data.Action != "reset" {
		return errors.ErrAccessDenied
	}

	err = repository.ChangePassword(data.UserID, params.Password, false)
	if err != nil {
		return err
	}

	return errors.ErrSuccess
}

// ChangePassword godoc
// @Summary Change password
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param password body string true "Password"
// @Success 200 {object} response.Message
// @Router /users/auth/password [put]
func ChangePassword(c echo.Context) (err error) {

	params := new(ChangePasswordShcema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	user := request.AuthenticatedUser(c)

	if err = repository.ChangePassword(user.ID, params.Password, false); err != nil {
		return err
	}

	return errors.ErrSuccess
}

// ChangeUsername godoc
// @Summary Update username
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param username body string true "Username"
// @Success 200 {object} response.Message
// @Router /users/auth/username [put]
func ChangeUsername(c echo.Context) (err error) {

	params := new(ChangeUsernameShcema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	user := request.AuthenticatedUser(c)

	if err = repository.ChangeUsername(user.ID, params.Username, false); err != nil {
		return err
	}

	return errors.ErrSuccess
}

// PutAvatar godoc
// @Summary Upload avatar
// @Tags user
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Message
// @Router /users/auth/avatar [put]
func PutAvatar(c echo.Context) (err error) {

	user := request.AuthenticatedUser(c)

	fh, err := c.FormFile("picture")
	if err != nil {
		return errors.ErrInvalidParams
	}

	f, err := fh.Open()
	if err != nil {
		return errors.ErrInternal
	}
	defer f.Close()

	if err = validation.ValidatePicture(fh.Size, f, config.AvatarPictureFormats, config.AvatarPictureMaxSize); err != nil {
		return err
	}

	// Reset file seeker
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return errors.ErrInternal
	}

	if err = storage.Put(user.ID, "avatar", f, fh.Size, fh.Header.Get("Content-Type")); err != nil {
		return err
	}

	return errors.ErrSuccess
}

// DeleteAvatar godoc
// @Summary Delete user's avatar
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Message
// @Router /users/auth/avatar [delete]
func DeleteAvatar(c echo.Context) (err error) {

	user := request.AuthenticatedUser(c)

	if err := storage.Delete(user.ID, "avatar"); err != nil {
		return err
	}

	return errors.ErrSuccess
}

// GetAccount godoc
// @Summary Get user account info
// @Tags user
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Message
// @Router /users/auth/mine [get]
func GetAccount(c echo.Context) (err error) {

	user := request.AuthenticatedUser(c)

	data, err := repository.GetAccountInfo(user.ID)
	if err != nil {
		return err
	}

	return r.CustomErrorJson(http.StatusOK, data, c)
}

// ChangeEmail godoc
// @Summary Change user account email
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param email body string true "Email"
// @Success 200 {object} response.Message
// @Router /users/auth/email [put]
func ChangeEmail(c echo.Context) (err error) {

	params := new(ChangeEmailShcema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	user := request.AuthenticatedUser(c)

	if err = repository.ChangeEmail(user.ID, params.Email, false); err != nil {
		return err
	}

	return errors.ErrSuccess
}

// GetUser godoc
// @Summary Get user with id
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "userID"
// @Success 200 {object} response.Message
// @Router /users/get/{id} [get]
func GetUser(c echo.Context) error {

	userID := c.Param("id")

	user, err := repository.GetUserByID(userID)
	if err != nil {
		return err
	}

	return r.CustomErrorJson(http.StatusOK, user, c)
}
