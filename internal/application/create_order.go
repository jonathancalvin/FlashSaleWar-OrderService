package application

import (
	"context"

	"github.com/jonathancalvin/golang-clean-architecture-template/internal/infrastructure/repository"
)

type CreateOrderUseCase struct {
	repo repository.Repository[any]
}

func NewCreateOrderUseCase(repo repository.Repository[any]) *CreateOrderUseCase {
	return &CreateOrderUseCase{repo: repo}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context) error {
	// business orchestration here
	return nil
}