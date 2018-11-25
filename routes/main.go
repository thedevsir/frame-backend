package routes

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/swaggo/echo-swagger"
	c "github.com/thedevsir/frame-backend/app/controller"
	"github.com/thedevsir/frame-backend/config"
	"github.com/thedevsir/frame-backend/middleware/auth"
	"github.com/thedevsir/frame-backend/middleware/objectId"
	"github.com/thedevsir/frame-backend/services/response"
)

func Composer() *echo.Echo {
	Route := echo.New()
	Route.HTTPErrorHandler = response.ErrorHandler
	endpoints := Route.Group("/endpoint")
	endpoints.Use(objectId.Check)
	{
		User := endpoints.Group("/users")
		{
			User.GET("/get/:id", c.GetUser).Name = "client get-user"
			User.POST("/signup", c.Signup).Name = "client new-user"
			User.POST("/signup/resend", c.Resend).Name = "client send-verification-email"
			User.POST("/signup/verification", c.Verification).Name = "client check-verification-token"
			User.POST("/signin", c.Signin).Name = "client let-user-in"
			User.POST("/signin/forgot", c.Forgot).Name = "client forgot-password"
			User.PUT("/signin/reset", c.Reset).Name = "client reset-password"
			{
				Auth := User.Group("/auth")
				Auth.Use(middleware.JWT([]byte(config.SigningKey)))
				Auth.Use(auth.Middleware)
				Auth.GET("/mine", c.GetAccount).Name = "client get-account"
				Auth.PUT("/username", c.ChangeUsername).Name = "client change-username"
				Auth.PUT("/email", c.ChangeEmail).Name = "client change-email"
				Auth.PUT("/password", c.ChangePassword).Name = "client change-password"
				Auth.PUT("/avatar", c.PutAvatar).Name = "client put-avatar"
				Auth.DELETE("/avatar", c.DeleteAvatar).Name = "client delete-avatar"
				Auth.GET("/sessions", c.Sessions).Name = "client get-sessions"
				Auth.DELETE("/signout", c.Signout).Name = "client delete-session"
			}
		}
		Admin := endpoints.Group("/admin")
		{
			Admin.POST("/signin", c.AdminSignin).Name = "admin let-admin-in"
			Auth := Admin.Group("/auth")
			Auth.Use(middleware.JWT([]byte(config.AdminSigningKey)))
			Auth.Use(auth.AdminMiddleware)
			User := Auth.Group("/users")
			{
				User.GET("/get/all", c.AdminGetAllUsers).Name = "admin get-users"
				User.GET("/get/:id", c.AdminGetUser).Name = "admin get-user"
				User.PUT("/status/:id", c.AdminChangeUserStatus).Name = "admin change-user-status"
				User.PUT("/email/:id", c.AdminChnageEmail).Name = "admin change-email"
				User.PUT("/username/:id", c.AdminChangeUsername).Name = "admin change-username"
				User.PUT("/password/:id", c.AdminChangePassword).Name = "admin change-password"
				User.PUT("/avatar/:id", c.AdminPutAvatar).Name = "admin put-avatar"
				User.DELETE("/avatar/:id", c.AdminDeleteAvatar).Name = "admin delete-avatar"
			}
			AdminManage := Auth.Group("/admin-manage")
			{
				AdminManage.GET("/get/all", c.GetAllAdmins).Name = "admin get-admins"
				AdminManage.GET("/get/:id", c.GetAdmin).Name = "admin get-admin"
				AdminManage.POST("/create", c.CreateAdmin).Name = "admin new-admin"
				AdminManage.PUT("/status/:id", c.ChangeAdminStatus).Name = "admin change-admin-status"
				AdminManage.PUT("/username/:id", c.ChangeAdminUsername).Name = "admin update-admin"
				AdminManage.PUT("/password/:id", c.ChangeAdminPassword).Name = "admin update-admin-password"
			}
			Auth.DELETE("/signout", c.AdminSignout).Name = "admin delete-session"
		}
	}
	if config.Mode == "DEV" {
		Route.GET("/swagger/*", echoSwagger.WrapHandler).Name = "docs-swagger"
	}
	return Route
}
