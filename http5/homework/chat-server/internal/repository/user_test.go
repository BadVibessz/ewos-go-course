package repository

import (
	"context"
	"math"
	"testing"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/model"

	inmemory "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/db/in-memory"
	sliceutils "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/pkg/utils/slice"
)

func initRepo(ctx context.Context) *UserRepoInMemDB {
	db, _ := inmemory.NewInMemDB(ctx, "")

	return NewInMemUserRepo(db)
}

func isEqualCreateModelToUser(createModel *model.User, user *model.User) bool {
	return createModel.Email == user.Email &&
		createModel.Username == user.Username &&
		createModel.HashedPassword == user.HashedPassword
}

func TestUserCreatedPositive(t *testing.T) {
	ctx := context.Background()
	repo := initRepo(ctx)

	toCreate := model.User{
		Email:          "test@mail.com",
		Username:       "test",
		HashedPassword: "NoHash",
	}

	created, err := repo.AddUser(ctx, toCreate)
	if err != nil {
		t.Fatal()
	}

	if !isEqualCreateModelToUser(&toCreate, created) {
		t.Fatal()
	}

	_, err = repo.DeleteUser(ctx, created.ID)
	if err != nil {
		t.Fatal("cannot delete user")
	}
}

func TestGetAllUsersPositive(t *testing.T) {
	ctx := context.Background()
	repo := initRepo(ctx)

	toCreate1 := model.User{
		Email:          "test@mail.com",
		Username:       "test",
		HashedPassword: "NoHash",
	}

	toCreate2 := model.User{
		Email:          "test@mail.com2",
		Username:       "test2",
		HashedPassword: "NoHash2",
	}

	created1, err := repo.AddUser(ctx, toCreate1)
	if err != nil {
		t.Fatalf("cannot add user")
	}

	created2, err := repo.AddUser(ctx, toCreate2)
	if err != nil {
		t.Fatalf("cannot add user")
	}

	got := repo.GetAllUsers(ctx, 0, math.MaxInt64)
	if len(got) == 0 {
		t.Fatal()
	}

	if !sliceutils.ContainsValue(got, *created1) || !sliceutils.ContainsValue(got, *created2) {
		t.Fatal()
	}

	_, err = repo.DeleteUser(ctx, created1.ID)
	if err != nil {
		t.Fatal("cannot delete user")
	}

	_, err = repo.DeleteUser(ctx, created2.ID)
	if err != nil {
		t.Fatal("cannot delete user")
	}
}

func TestGetUserByIdPositive(t *testing.T) {
	ctx := context.Background()
	repo := initRepo(ctx)

	toCreate := model.User{
		Email:          "test@mail.com",
		Username:       "test",
		HashedPassword: "NoHash",
	}

	created, err := repo.AddUser(ctx, toCreate)
	if err != nil {
		t.Fatalf("cannot add user")
	}

	got, err := repo.GetUserByID(ctx, created.ID)
	if err != nil {
		t.Fatal()
	}

	if *got != *created {
		t.Fatalf("expected user not equals to actual")
	}

	_, err = repo.DeleteUser(ctx, created.ID)
	if err != nil {
		t.Fatal("cannot delete user")
	}
}

func TestGetUserByEmailPositive(t *testing.T) {
	ctx := context.Background()
	repo := initRepo(ctx)

	toCreate := model.User{
		Email:          "test@mail.com",
		Username:       "test",
		HashedPassword: "NoHash",
	}

	created, err := repo.AddUser(ctx, toCreate)
	if err != nil {
		t.Fatalf("cannot add user")
	}

	got, err := repo.GetUserByEmail(ctx, created.Email)
	if err != nil {
		t.Fatal()
	}

	if *got != *created {
		t.Fatalf("expected user not equals to actual")
	}

	_, err = repo.DeleteUser(ctx, created.ID)
	if err != nil {
		t.Fatal("cannot delete user")
	}
}

func TestGetUserByUsernamePositive(t *testing.T) {
	ctx := context.Background()
	repo := initRepo(ctx)

	toCreate := model.User{
		Email:          "test@mail.com",
		Username:       "test",
		HashedPassword: "NoHash",
	}

	created, err := repo.AddUser(ctx, toCreate)
	if err != nil {
		t.Fatalf("cannot add user")
	}

	got, err := repo.GetUserByUsername(ctx, created.Username)
	if err != nil {
		t.Fatal()
	}

	if *got != *created {
		t.Fatalf("expected user not equals to actual")
	}

	_, err = repo.DeleteUser(ctx, created.ID)
	if err != nil {
		t.Fatal("cannot delete user")
	}
}
