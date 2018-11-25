package controller

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/thedevsir/frame-backend/app/repository"
	"github.com/thedevsir/frame-backend/services/errors"
	"github.com/thedevsir/frame-backend/services/paginate"
	"github.com/thedevsir/frame-backend/services/request"
	r "github.com/thedevsir/frame-backend/services/response"
)

type (
	CreateAdminSchema struct {
		Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
		Password string `json:"password" validate:"required,min=8,max=50"`
	}
	ChangeAdminUsernameSchema struct {
		Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	}
	ChangeAdminPasswordSchema struct {
		Password string `json:"password" validate:"required,min=8,max=50"`
	}
	ChangeAdminStatusSchema struct {
		IsActive bool `json:"isActive"`
	}
)

// GetAllAdmins godoc
// @Summary Get all admins
// @Tags adminManage
// @Accept json
// @Produce json
// @Security AdminApiKeyAuth
// @Param page query number false "Page"
// @Param limit query number false "Limit"
// @Success 200 {object} response.Message
// @Router /admin/auth/admin-manage/get/all [get]
func GetAllAdmins(c echo.Context) (err error) {

	if !request.IsRootAdmin(c) {
		return errors.ErrAccessDenied
	}

	page, limit := paginate.HandleQueries(c)
	admins, err := repository.GetAdmins(page, limit)
	if err != nil {
		return err
	}

	return r.CustomErrorJson(http.StatusOK, admins, c)
}

// CreateAdmin godoc
// @Summary Create a new Admin
// @Tags adminManage
// @Accept json
// @Produce json
// @Security AdminApiKeyAuth
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 201 {object} response.Message
// @Router /admin/auth/admin-manage/create [post]
func CreateAdmin(c echo.Context) (err error) {

	if !request.IsRootAdmin(c) {
		return errors.ErrAccessDenied
	}

	params := new(CreateAdminSchema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	_, err = repository.CreateAdmin(params.Username, params.Password)
	if err != nil {
		return err
	}

	return errors.ErrCreated
}

// GetAdmin godoc
// @Summary Get admin by id
// @Tags adminManage
// @Accept json
// @Produce json
// @Security AdminApiKeyAuth
// @Param id path string true "adminID"
// @Success 200 {object} response.Message
// @Router /admin/auth/admin-manage/get/{id} [get]
func GetAdmin(c echo.Context) (err error) {

	if !request.IsRootAdmin(c) {
		return errors.ErrAccessDenied
	}

	adminID := c.Param("id")

	admin, err := repository.GetAdminByID(adminID)
	if err != nil {
		return err
	}

	return r.CustomErrorJson(http.StatusOK, admin, c)
}

// ChangeAdminStatus godoc
// @Summary Set admin status
// @Tags adminManage
// @Accept json
// @Produce json
// @Security AdminApiKeyAuth
// @Param id path string true "adminID"
// @Param isActive body bool true "IsActive"
// @Success 200 {object} response.Message
// @Router /admin/auth/admin-manage/status/{id} [put]
func ChangeAdminStatus(c echo.Context) (err error) {

	if !request.IsRootAdmin(c) {
		return errors.ErrAccessDenied
	}

	adminID := c.Param("id")
	if adminID == "000000000000000000000000" {
		return errors.ErrAccessDenied
	}

	params := new(ChangeAdminStatusSchema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	err = repository.ChangeAdminStatus(adminID, params.IsActive)
	if err != nil {
		return err
	}

	return errors.ErrSuccess
}

// ChangeAdminUsername godoc
// @Summary Update username of admin
// @Tags adminManage
// @Accept json
// @Produce json
// @Security AdminApiKeyAuth
// @Param id path string true "adminID"
// @Param username body string true "Username"
// @Success 200 {object} response.Message
// @Router /admin/auth/admin-manage/username/{id} [put]
func ChangeAdminUsername(c echo.Context) (err error) {

	if !request.IsRootAdmin(c) {
		return errors.ErrAccessDenied
	}

	adminID := c.Param("id")
	if adminID == "000000000000000000000000" {
		return errors.ErrAccessDenied
	}

	params := new(ChangeAdminUsernameSchema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	err = repository.AdminChangeUsername(adminID, params.Username)
	if err != nil {
		return err
	}

	return errors.ErrSuccess
}

// ChangeAdminPassword godoc
// @Summary Update password of admin
// @Tags adminManage
// @Accept json
// @Produce json
// @Security AdminApiKeyAuth
// @Param id path string true "adminID"
// @Param password body string true "Password"
// @Success 200 {object} response.Message
// @Router /admin/auth/admin-manage/password/{id} [put]
func ChangeAdminPassword(c echo.Context) (err error) {

	if !request.IsRootAdmin(c) {
		return errors.ErrAccessDenied
	}

	adminID := c.Param("id")

	params := new(ChangeAdminPasswordSchema)
	if err = request.GetInputs(c, params); err != nil {
		return err
	}

	err = repository.AdminChangePassword(adminID, params.Password)
	if err != nil {
		return err
	}

	return errors.ErrSuccess
}
