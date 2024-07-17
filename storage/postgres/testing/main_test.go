package testing

import (
	pb "auth/genproto/AuthService"
	"auth/models"
	"auth/storage/postgres"
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"
)

func Connection() *sql.DB {
	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func TestRegisterAuth(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	authUser := postgres.NewAuthUser(dbpool)
	req := &pb.RequestRegister{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "password",
		FullName: "Test User",
		UserType: "admin",
	}

	err := authUser.RegisterAuth(req)
	if err != nil {
		t.Fatal(err)
	}

	var count int
	err = dbpool.QueryRow(`SELECT COUNT(*) FROM users WHERE email=$1`, req.Email).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	authUser := postgres.NewAuthUser(dbpool)
	email := "testuser@example.com"

	userinfo, err := authUser.GetUserByEmail(email)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(userinfo)
}

func TestStoreRefreshToken(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	authUser := postgres.NewAuthUser(dbpool)
	refreshToken := &models.RefreshToken{
		UserId:    "b3e03a4d-e0bf-4458-9703-1d63058dbd6f",
		Token:     "refresh-token",
		ExpiresAt: 1721286508,
	}

	err := authUser.StoreRefreshToken(refreshToken)
	if err != nil {
		t.Fatal(err)
	}

	var count int
	err = dbpool.QueryRow(`SELECT COUNT(*) FROM refresh_token WHERE token=$1`, refreshToken.Token).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(count)
}

func TestUpdatePassword(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	authUser := postgres.NewAuthUser(dbpool)
	req := &pb.PasswordRequest{
		Email:    "testuser@example.com",
		Password: "newpassword",
	}

	err := authUser.UpdatePassword(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	var passwordHash string
	err = dbpool.QueryRow(`SELECT password_hash FROM users WHERE email=$1`, req.Email).Scan(&passwordHash)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(passwordHash)
}

func TestGetAllUser(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	authUser := postgres.NewAuthUser(dbpool)

	res, err := authUser.GetAllUser(&pb.Void{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestUpdateUser(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	authUser := postgres.NewAuthUser(dbpool)
	req := &pb.UpdateRequest{
		Id:       "b3e03a4d-e0bf-4458-9703-1d63058dbd6f",
		Username: "updateduser",
		Email:    "updateduser@example.com",
		Password: "updatedpassword",
		Fullname: "Updated User",
		Usertype: "user",
	}

	res, err := authUser.UpdateUser(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("User updated successfully", res.Message)

	var username, email string
	err = dbpool.QueryRow(`SELECT username, email FROM users WHERE id=$1`, req.Id).Scan(&username, &email)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(t, req.Username, username)
}

func TestUpdateUserType(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	authUser := postgres.NewAuthUser(dbpool)
	req := &pb.TypeUserRequest{
		Id:       "b3e03a4d-e0bf-4458-9703-1d63058dbd6f",
		Usertype: "newusertype",
	}

	res, err := authUser.UpdateUserType(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("User's user type has been updated", res.Message)

	var userType string
	err = dbpool.QueryRow(`SELECT user_type FROM users WHERE id=$1`, req.Id).Scan(&userType)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(req.Usertype, userType)
}

func TestGetUsers(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	authUser := postgres.NewAuthUser(dbpool)
	req := &pb.GetUsersRequest{
		Page:  1,
		Limit: 10,
	}

	res, err := authUser.GetUsers(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(t, res)
}

func TestDeleteUser(t *testing.T) {
	dbpool := Connection()
	defer dbpool.Close()

	authUser := postgres.NewAuthUser(dbpool)
	req := &pb.DeleteUserRequest{
		Id: "b3e03a4d-e0bf-4458-9703-1d63058dbd6f",
	}

	res, err := authUser.DeleteUser(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("The user is effectively deleted", res.Message)

	var deletedAt time.Time
	err = dbpool.QueryRow(`SELECT deleted_at FROM users WHERE id=$1`, req.Id).Scan(&deletedAt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(deletedAt)
}

