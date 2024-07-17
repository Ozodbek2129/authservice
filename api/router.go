package api

import (
	"auth/api/handler"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"

	_ "auth/api/docs"
)

// @title Auth API
// @version 1.0
// @description This is an authentication service API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:50052
// @BasePath /
func NewRouter(db *sql.DB, redisClient *redis.Client) *gin.Engine {
	router := gin.Default()

	h := handler.NewHanler(db, redisClient)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/registerauth", h.RegisterAuth)
	router.POST("/loginauth", h.LoginAuth)
	router.POST("/passwordrecovery", h.Passwordrecovery)
	router.POST("/verifycoderesetpassword", h.VerifyCodeAndResetPassword)
	router.POST("/updatetoken", h.UpdateToken)
	router.POST("/canceltoken", h.CancelToken)

	return router
}
