package sqlite

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"hip-hop-geek/internal/models"
)

const (
	addUserStmt = `
    INSERT INTO users(id, username, today_subscribe,
    releases_message_id, releases_page_count, today_releases_message_id, today_releases_page_count)
    VALUES (?, ?, ?, 0, 0, 0, 0);
    `

	getUserByUsernameQuery = `
    SELECT id, username, today_subscribe, releases_message_id,
    releases_page_count, today_releases_message_id, today_releases_page_count
    FROM users
    WHERE username=?;
    `

	setTodaySubscribeStmt = `
    UPDATE users
    SET today_subscribe = ?
    WHERE id = ?;
    `

	setReleasesMessageIdStmt = `
    UPDATE users
    SET releases_message_id = ?,
        releases_page_count = ?
    WHERE id = ?;
    `

	setTodayReleasesMessageIdStmt = `
    UPDATE users
    SET today_releases_message_id = ?,
        today_releases_page_count = ?
    WHERE id = ?;
    `

	getAllSubscribersQuery = `
    SELECT id, username, today_subscribe, releases_message_id,
    releases_page_count, today_releases_message_id, today_releases_page_count
    FROM users
    WHERE today_subscribe = true;
    `
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrFewUsersFound     = errors.New(
		"getting user by username returns few values, but expecting one",
	)
	ErrUserNotFound = errors.New("user not found")
)

type UserSqlite struct {
	Id                     int64  `db:"id"`
	Username               string `db:"username"`
	IsTodaySubscribe       bool   `db:"today_subscribe"`
	ReleasesMessageId      int64  `db:"releases_message_id"`
	ReleasesPageCount      int    `db:"releases_page_count"`
	TodayReleasesMessageId int64  `db:"today_releases_message_id"`
	TodayReleasesPageCount int    `db:"today_releases_page_count"`
}

type UsersSqliteRepo struct {
	DB *sqlx.DB
}

func NewUserSqliteRepo(db *sqlx.DB) *UsersSqliteRepo {
	return &UsersSqliteRepo{db}
}

func (u *UsersSqliteRepo) Close() {
	u.DB.Close()
}

func (u *UsersSqliteRepo) AddUser(user models.User) error {
	_, err := u.DB.Exec(addUserStmt, user.Id, user.Username, user.IsTodaySubscribe)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed:") {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("db error add user: %w", err)
	}

	return nil
}

func (u *UsersSqliteRepo) GetUserByUsername(username string) (*models.User, error) {
	var user []UserSqlite
	err := u.DB.Select(&user, getUserByUsernameQuery, username)
	if err != nil {
		return nil, fmt.Errorf("error while getting user by username: %w", err)
	}

	if len(user) > 1 {
		return nil, ErrFewUsersFound
	} else if len(user) == 0 {
		return nil, ErrUserNotFound
	}

	return &models.User{
		Id:                     user[0].Id,
		Username:               user[0].Username,
		IsTodaySubscribe:       user[0].IsTodaySubscribe,
		ReleasesMessageId:      user[0].ReleasesMessageId,
		ReleasesPageCount:      user[0].ReleasesPageCount,
		TodayReleasesMessageId: user[0].TodayReleasesMessageId,
		TodayReleasesPageCount: user[0].TodayReleasesPageCount,
	}, nil
}

func (u *UsersSqliteRepo) SetTodaySubscribe(userId int64, isSubscribe bool) error {
	_, err := u.DB.Exec(setTodaySubscribeStmt, isSubscribe, userId)
	if err != nil {
		return fmt.Errorf("error while trying to set subscribe: %w", err)
	}

	return nil
}

func (u *UsersSqliteRepo) GetAllSubscribers() ([]*models.User, error) {
	var users []UserSqlite
	err := u.DB.Select(&users, getAllSubscribersQuery)
	if err != nil {
		return nil, fmt.Errorf("error while getting all subscribers: %w", err)
	}

	if len(users) == 0 {
		return nil, ErrUserNotFound
	}

	usersResult := make([]*models.User, 0, len(users))
	for _, user := range users {
		usersResult = append(usersResult, &models.User{
			Id:                     user.Id,
			Username:               user.Username,
			IsTodaySubscribe:       user.IsTodaySubscribe,
			ReleasesMessageId:      user.ReleasesMessageId,
			ReleasesPageCount:      user.ReleasesPageCount,
			TodayReleasesMessageId: user.TodayReleasesMessageId,
			TodayReleasesPageCount: user.TodayReleasesPageCount,
		})
	}

	return usersResult, nil
}

func (u *UsersSqliteRepo) SetUserState(
	userId int64,
	messageType,
	messageId,
	pageCount int,
) error {
	var stmt string
	switch messageType {
	case models.ReleasesMessage:
		stmt = setReleasesMessageIdStmt
	case models.TodayReleasesMessage:
		stmt = setTodayReleasesMessageIdStmt
	}
	if stmt == "" {
		return fmt.Errorf(
			"error while switch case message type in users repo, messageType=%d",
			messageType,
		)
	}

	_, err := u.DB.Exec(stmt, messageId, pageCount, userId)
	if err != nil {
		return fmt.Errorf("error while trying to set messageId and pageCount in user: %w", err)
	}

	return nil
}
