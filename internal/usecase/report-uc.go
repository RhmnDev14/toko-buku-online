package usecase

import (
	"context"
	"toko_buku_online/internal/entity"
	"toko_buku_online/internal/logger"
	"toko_buku_online/internal/repository"
)

type ReportUc interface {
	GetSalesReport(ctx context.Context) (entity.Sales, error)
	GetBestSeller(ctx context.Context) ([]entity.BestBook, error)
	GetPriceReport(ctx context.Context) (entity.PriceBook, error)
}

type reportUc struct {
	log  logger.Logger
	repo repository.ReportRepo
}

func NewReportUc(log logger.Logger, repo repository.ReportRepo) ReportUc {
	return &reportUc{
		log:  log,
		repo: repo,
	}
}

func (u *reportUc) GetSalesReport(ctx context.Context) (entity.Sales, error) {
	u.log.Info("get sales report in uc", ctx)
	return u.repo.GetSalesReport(ctx)
}

func (u *reportUc) GetBestSeller(ctx context.Context) ([]entity.BestBook, error) {
	u.log.Info("get best seller in uc", ctx)
	return u.repo.GetBestSeller(ctx)
}

func (u *reportUc) GetPriceReport(ctx context.Context) (entity.PriceBook, error) {
	u.log.Info("get price report in uc", ctx)
	return u.repo.GetPriceReport(ctx)
}
