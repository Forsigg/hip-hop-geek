package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"hip-hop-geek/internal/models"
)

func TestUserSqliteAddUser(t *testing.T) {
	t.Run("check user creates and getting correct", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		repo := NewUserSqliteRepo(db)
		userModel := models.User{
			Id:               1,
			Username:         "forsigg",
			IsTodaySubscribe: false,
		}
		err := repo.AddUser(userModel)
		assert.NoError(t, err)

		userFromDb, err := repo.GetUserByUsername(userModel.Username)
		if !assert.NoError(t, err) {
			t.Fatal()
		}

		assert.Equal(t, &userModel, userFromDb)
	})

	t.Run("check correct error when add not unique username", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		repo := NewUserSqliteRepo(db)
		user := models.User{
			Id:               1,
			Username:         "forsigg",
			IsTodaySubscribe: false,
		}
		repo.AddUser(user)
		err := repo.AddUser(user)
		assert.ErrorIs(t, err, ErrUserAlreadyExists)
	})

	t.Run("if user not found by username", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		repo := NewUserSqliteRepo(db)
		user := models.User{
			Id:               1,
			Username:         "forsigg",
			IsTodaySubscribe: false,
		}
		_, err := repo.GetUserByUsername(user.Username)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("check set subscribe work correct", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		repo := NewUserSqliteRepo(db)
		user := &models.User{
			Id:               1,
			Username:         "forsigg",
			IsTodaySubscribe: false,
		}
		repo.AddUser(*user)

		err := repo.SetTodaySubscribe(user.Id, true)
		assert.NoError(t, err)

		userDb, _ := repo.GetUserByUsername(user.Username)
		assert.Equal(t, userDb.IsTodaySubscribe, true)
	})
}

func TestGetAllSubscribers(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		repo := NewUserSqliteRepo(db)
		user := models.User{
			Id:               1,
			Username:         "forsigg",
			IsTodaySubscribe: true,
		}
		repo.AddUser(user)

		users, err := repo.GetAllSubscribers()
		if err != nil {
			t.Fatal(err)
		}

		assert.NoError(t, err)
		assert.Equal(t, 1, len(users))
		assert.Equal(t, user, *users[0])
	})

	t.Run("if subscribers not found", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		repo := NewUserSqliteRepo(db)

		users, err := repo.GetAllSubscribers()
		assert.ErrorIs(t, err, ErrUserNotFound)
		assert.Nil(t, users)
	})
}

func TestSetUserState(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		repo := NewUserSqliteRepo(db)
		user := models.User{
			Id:               1,
			Username:         "forsigg",
			IsTodaySubscribe: true,
		}
		repo.AddUser(user)

		err := repo.SetUserState(
			user.Id,
			models.TodayReleasesMessage,
			1,
			5,
		)
		assert.NoError(t, err)

		userFromDb, _ := repo.GetUserByUsername(user.Username)
		assert.Equal(t, userFromDb.TodayReleasesMessageId, int64(1))
		assert.Equal(t, userFromDb.TodayReleasesPageCount, 5)
	})
}
