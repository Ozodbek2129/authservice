package handler

import (
	"auth/api/email"
	"auth/api/token"
	pb "auth/genproto/AuthService"
	"auth/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)


// RegisterAuth godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param  body  body  AuthService.RequestRegister  true  "Register Request"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /registerauth [post]
func (h *Handler) RegisterAuth(c *gin.Context) {
	req := pb.RequestRegister{}
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Password = string(hashpassword)

	err = h.AuthUser.RegisterAuth(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "user registered successfully"})
}

// LoginAuth godoc
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param  body  body  AuthService.RequestLogin  true  "Login Request"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /loginauth [post]
func (h *Handler) LoginAuth(c *gin.Context) {
	req := pb.RequestLogin{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	user, err := h.AuthUser.GetUserByEmail(req.Email)
	if err != nil {
		fmt.Println("124")
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	
	fmt.Println(user.Password)
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Password comparison failed: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := token.GenerateJWT(&pb.User{
		Id:       user.Id,
		Username: user.UserName,
		Email:    req.Email,
		Fullname: user.FullName,
		Usertype: user.UserType,
	})

	err = h.AuthUser.StoreRefreshToken(&models.RefreshToken{
		UserId:    user.Id,
		Token:     token.Refreshtoken,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	if err != nil {
		fmt.Println("126")
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusAccepted, token)
}

// Passwordrecovery godoc
// @Summary Recover password
// @Description Send password recovery email
// @Tags auth
// @Accept  json
// @Produce  json
// @Param  body  body  AuthService.PasswordRequest  true  "Password Recovery Request"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /passwordrecovery [post]
func (h *Handler) Passwordrecovery(c *gin.Context) {
	req := pb.PasswordRequest{}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	ctx := context.Background()

	err = h.Redis.Set(ctx, req.Email, code, time.Minute*8).Err()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	email.SendCode(req.Email, code)
	c.JSON(http.StatusAccepted, gin.H{"message": "Password recovery email sent"})
}

// VerifyCodeAndResetPassword godoc
// @Summary Verify code and reset password
// @Description Verify the recovery code and reset the password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param  body  body  models.Request  true  "Verify Code and Reset Password Request"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /verifycoderesetpassword [post]
func (h Handler) VerifyCodeAndResetPassword(c *gin.Context) {
	req := models.Request{}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	ctx := context.Background()
	storedcode, err := h.Redis.Get(ctx, req.Email).Result()
	if err == redis.Nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired code"})
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if storedcode != req.Code {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	err = h.AuthUser.UpdatePassword(context.Background(), &pb.PasswordRequest{
		Email:    req.Email,
		Password: string(hashpassword),
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	h.Redis.Del(ctx, req.Email)
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Password Update successfully",
	})
}

// UpdateToken godoc
// @Summary Update access token
// @Description Update the access token using the refresh token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param  body  body  AuthService.RefreshToken  true  "Update Token Request"
// @Success 202 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /updatetoken [post]
func (h *Handler) UpdateToken(c *gin.Context) {
	req := pb.RefreshToken{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newtoken, err := token.RefreshJWT(req.Refreshtoken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.AuthUser.StoreRefreshToken(&models.RefreshToken{
		UserId:    newtoken.Userid,
		Token:     newtoken.Refreshtoken,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, newtoken)
}


// CancelToken godoc
// @Summary Cancel refresh token
// @Description Cancel the refresh token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param  body  body  AuthService.RefreshToken  true  "Cancel Token Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /canceltoken [post]
func (h *Handler) CancelToken(c *gin.Context) {
	req := pb.RefreshToken{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	resp, err := h.AuthUser.DeleteRefreshToken(req.Refreshtoken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
