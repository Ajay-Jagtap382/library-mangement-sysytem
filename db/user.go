package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

const (
	createUserQuery = `INSERT INTO user (
		id,
		first_name,
		last_name,
		mobile_num ,
		email,
		password,
		gender,
		role
	)
        VALUES(?,?,?,?,?,?,?,?)`

	idexist = `SELECT COUNT(*) FROM user WHERE user.id = ?`

	listUserQuery        = `SELECT * FROM user`
	findUserByIDQuery    = `SELECT * FROM user WHERE id = ?`
	findUserByEmailQuery = `SELECT * FROM user WHERE email = ?`
	deleteUserByIDQuery  = `DELETE FROM user WHERE id = ?`
	updateUserQuery      = `UPDATE user SET first_name = ?, last_name = ? WHERE id = ?`
	updatePasswordQuery  = `UPDATE user SET password=? where id=?`
)

type User struct {
	ID         string `db:"id"`
	First_Name string `db:"first_name"`
	Last_Name  string `db:"last_name"`
	Mobile_Num string `db:"mobile_num"`
	Email      string `db:"email"`
	Password   string `db:"password"`
	Gender     string `db:"gender"`
	Role       string `db:"role"`
}

func (s *store) CreateUser(ctx context.Context, user *User) (err error) {

	return Transact(ctx, s.db, &sql.TxOptions{}, func(ctx context.Context) error {
		_, err = s.db.Exec(
			createUserQuery,
			user.ID,
			user.First_Name,
			user.Last_Name,
			user.Mobile_Num,
			user.Email,
			user.Password,
			strings.ToLower(user.Gender),
			strings.ToLower(user.Role),
		)
		return err
	})
}

func (s *store) ListUsers(ctx context.Context) (users []User, err error) {
	err = WithDefaultTimeout(ctx, func(ctx context.Context) error {
		return s.db.SelectContext(ctx, &users, listUserQuery)
	})
	if err == sql.ErrNoRows {
		return users, ErrUserNotExist
	}
	return
}

func (s *store) FindUserByID(ctx context.Context, id string) (user User, err error) {
	err = WithDefaultTimeout(ctx, func(ctx context.Context) error {
		return s.db.GetContext(ctx, &user, findUserByIDQuery, id)
	})
	if err == sql.ErrNoRows {
		return user, ErrUserNotExist
	}
	return
}

func (s *store) FindUserByEmail(ctx context.Context, email string) (user User, err error) {
	err = WithDefaultTimeout(ctx, func(ctx context.Context) error {
		return s.db.GetContext(ctx, &user, findUserByEmailQuery, email)
	})
	if err == sql.ErrNoRows {
		return user, ErrUserNotExist
	}
	return
}

func (s *store) DeleteUserByID(ctx context.Context, id string) (err error) {
	return Transact(ctx, s.db, &sql.TxOptions{}, func(ctx context.Context) error {
		res, err := s.db.Exec(deleteUserByIDQuery, id)
		cnt, err := res.RowsAffected()
		if cnt == 0 {
			return ErrUserNotExist
		}
		if err != nil {
			return err
		}
		return err
	})
}

func (s *store) UpdateUser(ctx context.Context, user *User) (err error) {

	flag := 0

	s.db.GetContext(ctx, &flag, idexist, user.ID)
	fmt.Println(flag)

	if flag == 0 {
		return ErrIDNotExist
	} else {
		return Transact(ctx, s.db, &sql.TxOptions{}, func(ctx context.Context) error {
			_, err = s.db.Exec(
				updateUserQuery,
				user.First_Name,
				user.Last_Name,
				user.ID,
			)
			return err
		})
	}
}

func (s *store) UpdatePassword(ctx context.Context, user *User) (err error) {

	flag := 0

	s.db.GetContext(ctx, &flag, idexist, user.ID)

	if flag == 0 {
		return ErrIDNotExist
	} else {

		return Transact(ctx, s.db, &sql.TxOptions{}, func(ctx context.Context) error {
			_, err = s.db.Exec(
				updatePasswordQuery,
				user.Password,
				user.ID,
			)
			return err
		})
	}
}
