package users

import (
	"bank-api/internal/domain"
	"context"
)

type ReadUserUseCase struct {
	userRepository domain.UserRepository
}

func NewReadUserUseCase(userRepository domain.UserRepository) *ReadUserUseCase {
	return &ReadUserUseCase{
		userRepository: userRepository,
	}
}

type ReadUserCommand struct {
	Login string
}

func (useCase ReadUserUseCase) ReadUserHandler(ctx context.Context, command *ReadUserCommand) (*domain.User, error) {
	user, err := useCase.userRepository.FindByLogin(ctx, command.Login)
	return user, err
}
