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
