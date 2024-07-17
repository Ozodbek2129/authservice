package postgres

import (
	pb "auth/genproto/AuthService"
	"auth/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type AuthUser struct {
	db *sql.DB
}

func NewAuthUser(db *sql.DB) *AuthUser {
	return &AuthUser{db: db}
}

func (u *AuthUser) RegisterAuth(req *pb.RequestRegister) error {
	id := uuid.NewString()
	_, err := u.db.Exec(`
		insert into users(
			id,username, email, password_hash, full_name ,user_type
		)values(
			$1,$2,$3,$4,$5,$6)`, id, req.Username, req.Email, req.Password, req.FullName, req.UserType)

	if err != nil {
		log.Fatalf("User register error: %v", err)
		return err
	}

	return nil
}

func (u *AuthUser) GetUserByEmail(email string) (*models.UserInfo, error) {
	userinfo := models.UserInfo{}
	err := u.db.QueryRow(
		`select 
			id, username, password_hash, full_name, user_type
		from
			users
		where 
			email=$1 and 
			deleted_at is null`, email).Scan(
		&userinfo.Id, &userinfo.UserName, &userinfo.Password, &userinfo.FullName, &userinfo.UserType)

	if err != nil {
		return nil, err
	}

	return &userinfo, nil
}

func (u *AuthUser) StoreRefreshToken(req *models.RefreshToken) error {
	id := uuid.NewString()
	_, err := u.db.Exec(
		`insert into refresh_token(
			id, user_id, token, expires_at
		)values(
			$1,$2,$3,$4)`, id, req.UserId, req.Token, req.ExpiresAt)

	if err != nil {
		log.Fatalf("Error with inserting refresh_token: %v", err)
		return err
	}
	return nil
}

func (u *AuthUser) UpdatePassword(ctx context.Context, req *pb.PasswordRequest) error {
	query := `update users 
			set password_hash=$1
			where email=$2`
	_, err := u.db.ExecContext(ctx, query, req.Password, req.Email)
	if err != nil {
		log.Fatalf("Error resetting password: %v", err)
	}
	return nil
}

func (u AuthUser) DeleteRefreshToken(refreshtoken string) (*pb.Status, error) {
	query := `delete from refresh_token
	        where token=$1`

	_, err := u.db.Exec(query, refreshtoken)
	if err != nil {
		return nil, err
	}

	return &pb.Status{Message: "The token has been revoked"}, nil
}

func (u *AuthUser) GetAllUser(req *pb.Void) (*pb.GetAllResponse, error) {
	query := `select 
				id, username, email, password_hash, full_name, user_type
			from
				users`

	rows, err := u.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		var user pb.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Fullname, &user.Usertype); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pb.GetAllResponse{Users: users}, nil
}

func (u *AuthUser) UpdateUser(req *pb.UpdateRequest) (*pb.Status, error) {
	query := `UPDATE users SET 
				username = COALESCE(NULLIF($1, ''), username), 
				email = COALESCE(NULLIF($2, ''), email), 
				password_hash = COALESCE(NULLIF($3, ''), password_hash), 
				full_name = COALESCE(NULLIF($4, ''), full_name), 
				user_type = COALESCE(NULLIF($5, ''), user_type)
			WHERE id = $6`

	_, err := u.db.Exec(query, req.Username, req.Email, req.Password, req.Fullname, req.Usertype, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Status{
		Message: "User updated successfully",
	}, nil
}

func (u *AuthUser) UpdateUserType(req *pb.TypeUserRequest) (*pb.Status, error) {
	query := `update users set
				user_type=$1
			where id=$2`

	_, err := u.db.Exec(query, req.Usertype, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Status{Message: "User's user type has been updated"}, nil
}

func (u *AuthUser) GetUsers(req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	fmt.Print("sfndsdsndsodsndsnodsondsoosso")
	if req.Page < 1 {
		return nil, fmt.Errorf("invalid page number")
	}

	offset := (req.Page - 1) * req.Limit
	query := `select id, username, email, password_hash, full_name, user_type
              from users
              LIMIT $1 OFFSET $2`
	rows, err := u.db.Query(query, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		var user pb.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Fullname, &user.Usertype); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var total int
	err = u.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, err
	}

	return &pb.GetUsersResponse{
		Users: users,
		Total: int32(total),
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

func (u *AuthUser) DeleteUser(req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	query := `update users set
				deleted_at=$1
			where id=$2`

	newtime := time.Now()

	_, err := u.db.Exec(query, newtime, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteUserResponse{Message: "The user is effectively deleted"}, nil
}

func (u *AuthUser) IdCheck(req *pb.Id) (*pb.Response, error) {
	query := `select id from users`

	rows, err := u.db.Query(query)
	if err != nil {
		return &pb.Response{B: false}, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return &pb.Response{B: false}, err
		}
		if id == req.Id {
			return &pb.Response{B: true}, nil
		}
	}

	if err = rows.Err(); err != nil {
		return &pb.Response{B: false}, err
	}

	return &pb.Response{B: false}, nil
}

func (u *AuthUser) SearchName(req *pb.Id) (*pb.Name, error) {
	query := `select username from users where id=$1`

	var username string
	err := u.db.QueryRow(query, req.Id).Scan(&username)
	if err != nil {
		return nil, err
	}
	return &pb.Name{
		Name: username,
	}, nil
}
