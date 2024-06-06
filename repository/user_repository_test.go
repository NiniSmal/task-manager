package repository

import (
	"context"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"gitlab.com/nina8884807/task-manager/entity"
	"testing"
	"time"
)

func TestUserRepository_CreateUser(t *testing.T) {
	db, rds := DBConnection(t)
	ur := NewUserRepository(db, rds)
	photo := uuid.NewString()

	user := entity.User{
		Email:            uuid.NewString(),
		Password:         "123",
		CreatedAt:        time.Now().UTC().Round(time.Millisecond),
		Role:             "user",
		Verification:     true,
		VerificationCode: "41c37c27-291b-477d-af8d-c162e0fa3e98",
		Photo:            &photo,
	}
	ctx := context.Background()
	id, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)

	user.ID = id
	user.Password = ""
	user.VerificationCode = ""

	dbUser, err := ur.GetUserByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, user, dbUser)
}

func TestUserRepository_GetUserByID_Error(t *testing.T) {
	db, rds := DBConnection(t)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()
	_, err := ur.GetUserByID(ctx, 1234)
	require.Error(t, err)
}

func TestUserRepository_UserByLogin(t *testing.T) {
	db, rds := DBConnection(t)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()

	user := entity.User{
		Email:            uuid.New().String(),
		Password:         uuid.New().String(),
		CreatedAt:        time.Time{},
		Role:             "user",
		Verification:     true,
		VerificationCode: "123456",
	}
	id, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)

	dbUser, err := ur.UserByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.Equal(t, id, dbUser.ID)
}

func TestUserRepository_Verification(t *testing.T) {
	db, rds := DBConnection(t)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()
	user := entity.User{
		Email:            uuid.New().String(),
		Password:         uuid.New().String(),
		CreatedAt:        time.Time{},
		Role:             "user",
		Verification:     false,
		VerificationCode: "123",
	}
	_, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)

	userUp := entity.User{
		Email:            user.Email,
		Password:         user.Password,
		CreatedAt:        user.CreatedAt,
		Role:             "user",
		Verification:     true,
		VerificationCode: "123",
	}
	dbID, err := ur.VerifyByCode(ctx, userUp.VerificationCode)
	require.NoError(t, err)

	user.ID = dbID
	userDB, err := ur.GetUserByID(ctx, dbID)

	require.NoError(t, err)
	require.Equal(t, user.ID, userDB.ID)
	require.Equal(t, userDB.Verification, userUp.Verification)

}

func TestUserRepository_SaveSession(t *testing.T) {
	db, rds := DBConnection(t)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()
	session := uuid.New()
	createdAtSession := time.Now()
	user := entity.User{
		ID:       int64(1),
		Email:    uuid.New().String(),
		Password: uuid.New().String(),
		Role:     "user",
	}
	err := ur.SaveSession(ctx, session, user, createdAtSession)
	require.NoError(t, err)

	userS, err := ur.GetSession(ctx, session)
	require.NoError(t, err)
	require.Equal(t, user, userS)
}
func TestUserRepository_SaveSessionFromCache(t *testing.T) {
	db, rds := DBConnection(t)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()
	session := uuid.New()
	user := entity.User{
		ID:       int64(1),
		Email:    uuid.New().String(),
		Password: uuid.New().String(),
		Role:     "user",
	}
	err := ur.saveSessionToCache(ctx, session, user)
	require.NoError(t, err)

	userS, err := ur.getSessionFromCache(ctx, session)
	require.NoError(t, err)
	require.Equal(t, user, userS)
}

func TestUserRepository_GetSession_Error(t *testing.T) {
	db, rds := DBConnection(t)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()

	_, err := ur.GetSession(ctx, uuid.New())
	require.Error(t, err)
}

func TestUserRepository_GetSessionFromCache_Error(t *testing.T) {
	db, rds := DBConnection(t)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()

	_, err := ur.getSessionFromCache(ctx, uuid.New())
	require.Error(t, err)
}

func TestUserRepository_UpdateVerificationCode(t *testing.T) {
	db, rds := DBConnection(t)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()

	user := entity.User{
		Email:            uuid.New().String(),
		Password:         uuid.New().String(),
		CreatedAt:        time.Time{},
		Role:             "user",
		Verification:     false,
		VerificationCode: "123",
	}

	_, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)

	code := "456"
	user.VerificationCode = code
	err = ur.UpdateVerificationCode(ctx, user.ID, user.VerificationCode)
	require.NoError(t, err)
	require.Equal(t, code, user.VerificationCode)

}
func TestUserRepository_SavePhoto(t *testing.T) {
	db, rds := DBConnection(t)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()

	sessionID := uuid.New()
	ctx = context.WithValue(ctx, "session_id", sessionID)
	user := entity.User{
		Email:            uuid.New().String(),
		Password:         uuid.New().String(),
		CreatedAt:        time.Time{},
		Role:             "user",
		Verification:     false,
		VerificationCode: "123",
	}

	id, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)

	photo := uuid.NewString()

	err = ur.SavePhoto(ctx, photo, id)
	require.NoError(t, err)
}

func TestUserRepository_Users(t *testing.T) {
	db, rds := DBConnection(t)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()

	users := []entity.User{
		{Email: uuid.New().String(),
			Password:         uuid.New().String(),
			CreatedAt:        time.Time{},
			Role:             "user",
			Verification:     true,
			VerificationCode: "123"},
		{Email: uuid.New().String(),
			Password:         uuid.New().String(),
			CreatedAt:        time.Time{},
			Role:             "user",
			Verification:     true,
			VerificationCode: "1234"},
	}

	for i, user := range users {
		id, err := ur.CreateUser(ctx, user)
		require.NoError(t, err)
		users[i].ID = id
		users[i].Password = ""
	}

	usersDB, err := ur.Users(ctx)
	require.NoError(t, err)

	for _, user := range users {
		require.Contains(t, usersDB, user)
	}
}

func TestUserRepository_DeleteUser(t *testing.T) {
	db, rds := DBConnection(t)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()

	user := entity.User{Email: uuid.NewString()}
	userID, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	user.ID = userID

	err = ur.DeleteUser(ctx, user.ID)
	require.NoError(t, err)

	_, err = ur.GetUserByID(ctx, userID)
	require.Error(t, err)
}
