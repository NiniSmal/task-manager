package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gitlab.com/nina8884807/task-manager/entity"
	"testing"
	"time"
)

func UserConnection(t *testing.T) (*sql.DB, *redis.Client) {
	t.Helper() //помечает как вспомогвтельную
	// docker run -d -p 9000:5432 -e POSTGRES_PASSWORD=dev -e POSTGRES_DATABASE=postgres postgres
	db, err := sql.Open("postgres", "postgres://postgres:dev@localhost:9000/postgres?sslmode=disable")

	require.NoError(t, err)
	t.Cleanup(func() {
		err := db.Close()
		require.NoError(t, err)
	})
	err = db.Ping()
	require.NoError(t, err)

	ctx := context.Background()
	rds := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	t.Cleanup(func() {
		err := rds.Close()
		require.NoError(t, err)
	})

	_, err = rds.Ping(ctx).Result()
	require.NoError(t, err)
	return db, rds
}

func TestUserRepository_CreateUser(t *testing.T) {
	db, rds := UserConnection(t)
	ur := NewUserRepository(db, rds)
	user := entity.User{
		Email:            uuid.NewString(),
		Password:         "123",
		CreatedAt:        time.Now(),
		Role:             "user",
		Verification:     true,
		VerificationCode: "41c37c27-291b-477d-af8d-c162e0fa3e98",
	}
	ctx := context.Background()
	id, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)

	dbUser, err := ur.GetUserByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, user.ID, dbUser.ID)
	require.Equal(t, user.Email, dbUser.Email)
	require.Equal(t, user.Password, dbUser.Password)
	require.Equal(t, user.CreatedAt.Unix(), dbUser.CreatedAt.Unix())
	require.Equal(t, user.Role, dbUser.Role)
	require.Equal(t, user.Verification, dbUser.Verification)
	require.Equal(t, user.VerificationCode, dbUser.VerificationCode)
}

func TestUserRepository_GetUserByID_Error(t *testing.T) {
	db, rds := UserConnection(t)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()
	_, err := ur.GetUserByID(ctx, 1234)
	require.Error(t, err)
}

func TestUserRepository_UserByLogin(t *testing.T) {
	db, rds := UserConnection(t)
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

	dbUser, err := ur.UserByLogin(ctx, user.Email)
	require.NoError(t, err)
	require.Equal(t, id, dbUser.ID)
}

func TestUserRepository_Verification(t *testing.T) {
	db, rds := UserConnection(t)
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
	dbID, err := ur.Verification(ctx, userUp.VerificationCode, userUp.Verification)
	require.NoError(t, err)

	userDB, err := ur.GetUserByID(ctx, dbID)

	require.NoError(t, err)
	require.Equal(t, user.ID, userDB.ID)
	require.Equal(t, userDB.VerificationCode, userUp.VerificationCode)
	require.Equal(t, userDB.Verification, userUp.Verification)

}

func TestUserRepository_SaveSession(t *testing.T) {
	db, rds := UserConnection(t)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()
	session := uuid.New()
	user := entity.User{
		ID:       int64(1),
		Email:    uuid.New().String(),
		Password: uuid.New().String(),
		Role:     "user",
	}
	err := ur.SaveSession(ctx, session, user)
	require.NoError(t, err)

	userS, err := ur.GetSession(ctx, session)
	require.NoError(t, err)
	require.Equal(t, user, userS)
}

func TestUserRepository_GetSession_Error(t *testing.T) {
	db, rds := UserConnection(t)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()

	_, err := ur.GetSession(ctx, uuid.New())
	require.Error(t, err)
}
