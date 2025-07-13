package merchant

import (
	"context"
)

type MerchantUsecase interface {
	CreateMerchant(ctx context.Context, name, description string) (*Merchant, error)
}

type merchantUsecase struct {
	repo MerchantRepository
}

func NewMerchantUsecase(repo MerchantRepository) MerchantUsecase {
	return &merchantUsecase{
		repo: repo,
	}
}

func (u *merchantUsecase) CreateMerchant(ctx context.Context, name, description string) (*Merchant, error) {
	merchant := NewMerchant(name, description)
	err := u.repo.Save(ctx, merchant)
	if err != nil {
		return nil, err
	}
	return merchant, nil
}
