package users

import (
	"bank-api/internal/domain"
	"bank-api/internal/pkg"
	"context"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserUseCase struct {
	userRepository domain.UserRepository
}

func NewCreateUserUseCase(userRepository domain.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: userRepository,
	}
}

type CreateUserCommand struct {
	Login    string
	Password []byte
}

func (useCase *CreateUserUseCase) CreateUserHandler(ctx context.Context, command *CreateUserCommand) (*domain.User, error) {
	password, err := bcrypt.GenerateFromPassword(command.Password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := domain.NewUser(pkg.GenerateID(), command.Login, password)
	err = useCase.userRepository.Save(ctx, user)
	return user, err
}
