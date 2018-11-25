package controller

import (
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/thedevsir/frame-backend/app/repository"
	"github.com/thedevsir/frame-backend/config"
	"github.com/thedevsir/frame-backend/services/errors"
	"github.com/thedevsir/frame-backend/services/paginate"
	"github.com/thedevsir/frame-backend/services/request"
	r "github.com/thedevsir/frame-backend/services/response"
	"github.com/thedevsir/frame-backend/services/storage"
	"github.com/thedevsir/frame-backend/services/validation"
)

type (
	AdminChangeUsernameSchema struct {
		Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	}
	AdminChangePasswordSchema struct {
		Password string `json:"password" validate:"required,min=8,max=50"`
	}
	AdminChnageEmailSchema struct {
		Email string `json:"email" validate:"required,email"`
	}
	AdminChangeUserStatusSchema struct {
		IsActive bool `json:"isActive"`
	}
)

// AdminGetAllUsers godoc
// @Summary Get all users
// @Tags adminUser
// @Accept json
// @Produce json
// @Security AdminApiKeyAuth
// @Param page query int false "page of pgination"
// @Param limit query int false "limit of pgination"
// @Success 200 {object} response.Message
// @Router /admin/auth/users/get/all [get]
func AdminGetAllUsers(c echo.Context) (err error) {

	page, limit := paginate.HandleQueries(c)
	users, err := repository.GetUsers(page, limit)
	if err != nil {
		return err
	}

	return r.CustomErrorJson(http.StatusOK, users, c)
}

// AdminGetUser godoc
// @Summary Get user by id
// @Tags adminUser
// @Accept json
// @Produce json
// @Security AdminApiKeyAuth
// @Param id path string true "user ID"
// @Success 200 {object} response.Message
// @Router /admin/auth/users/get/{id} [get]
func AdminGetUser(c echo.Context) (err error) {

	userID := c.Param("id")
	user, err := repository.GetUserByIDFromAdmin(userID)
	if err != nil {
		return err
	}

	return r.CustomErrorJson(http.StatusOK, user, c)
}

// AdminChangeUsername godoc
// @Summary Set username
// @Tags adminUser
// @Accept json
// @Produce json
// @Security AdminApiKeyAuth
// @Param id path string true "user ID"
// @Param username body string true "Username"
// @Success 200 {object} response.Message
// @Router /admin/auth/users/username/{id} [put]
func AdminChangeUsername(c echo.Context) (err error) {

	userID := c.Param("id")

	params := new(AdminChangeUsernameSchema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	if err = repository.ChangeUsername(userID, params.Username, true); err != nil {
		return err
	}

	return errors.ErrSuccess
}

// AdminChangePassword godoc
// @Summary Set user password
// @Tags adminUser
// @Accept json
// @Produce json
// @Security AdminApiKeyAuth
// @Param id path string true "user ID"
// @Param password body string true "Password"
// @Success 200 {object} response.Message
// @Router /admin/auth/users/password/{id} [put]
func AdminChangePassword(c echo.Context) (err error) {

	userID := c.Param("id")

	params := new(AdminChangePasswordSchema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	if err = repository.ChangePassword(userID, params.Password, true); err != nil {
		return err
	}

	return errors.ErrSuccess
}

// AdminChnageEmail godoc
// @Summary Set user email
// @Tags adminUser
// @Accept json
// @Produce json
// @Security AdminApiKeyAuth
// @Param id path string true "user ID"
// @Param email body string true "Email"
// @Success 200 {object} response.Message
// @Router /admin/auth/users/email/{id} [put]
func AdminChnageEmail(c echo.Context) (err error) {

	userID := c.Param("id")

	params := new(AdminChnageEmailSchema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	if err = repository.ChangeEmail(userID, params.Email, true); err != nil {
		return err
	}

	return errors.ErrSuccess
}

// AdminPutAvatar godoc
// @Summary upload avatar
// @Tags adminUser
// @Accept multipart/form-data
// @Produce json
// @Security AdminApiKeyAuth
// @Param id path string true "user ID"
// @Success 200 {object} response.Message
// @Router /admin/auth/users/avatar/{id} [put]
func AdminPutAvatar(c echo.Context) (err error) {

	userID := c.Param("id")

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
		return err
	}

	if err = storage.Put(userID, "avatar", f, fh.Size, fh.Header.Get("Content-Type")); err != nil {
		return err
	}

	return errors.ErrSuccess
}

// AdminDeleteAvatar godoc
// @Summary delete user's avatar
// @Tags adminUser
// @Accept json
// @Produce json
// @Security AdminApiKeyAuth
// @Param id path string true "user ID"
// @Success 200 {object} response.Message
// @Router /admin/auth/users/avatar/{id} [delete]
func AdminDeleteAvatar(c echo.Context) (err error) {

	userID := c.Param("id")

	if err = storage.Delete(userID, "avatar"); err != nil {
		return err
	}

	return errors.ErrSuccess
}

// AdminChangeUserStatus godoc
// @Summary Set user status
// @Tags adminUser
// @Accept json
// @Produce json
// @Security AdminApiKeyAuth
// @Param id path string true "UserID"
// @Param isActive body bool true "IsActive"
// @Success 201 {object} response.Message
// @Router /admin/auth/users/status/{id} [put]
func AdminChangeUserStatus(c echo.Context) (err error) {

	userID := c.Param("id")

	params := new(AdminChangeUserStatusSchema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	if err = repository.ChangeUserStatus(userID, params.IsActive); err != nil {
		return err
	}

	return errors.ErrSuccess
}
