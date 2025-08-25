package handler

import (
	"context"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/logger"
	"toko_buku_online/internal/middleware"
	"toko_buku_online/internal/usecase"
	"toko_buku_online/toko_buku_online/api/gen/go/toko/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReportHandler struct {
	toko.UnimplementedReportServiceServer
	log        logger.Logger
	uc         usecase.ReportUc
	middleware middleware.AuthMiddleware
}

func NewReportHandler(log logger.Logger, uc usecase.ReportUc, middleware middleware.AuthMiddleware) *ReportHandler {
	return &ReportHandler{
		log:        log,
		uc:         uc,
		middleware: middleware,
	}
}

func (h *ReportHandler) GetSalesReport(ctx context.Context, req *toko.EmptyReport) (*toko.SalesReport, error) {
	h.log.Info("get sales report in handler", req)

	ctx, err := h.middleware.Require(ctx, constant.GET)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	report, err := h.uc.GetSalesReport(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &toko.SalesReport{
		Omset:            report.Omset,
		TotalBukuTerjual: int32(report.TotalBukuTerjual),
	}, nil
}
func (h *ReportHandler) GetBestSeller(ctx context.Context, req *toko.EmptyReport) (*toko.BestSellerReport, error) {
	h.log.Info("get sales report in handler", req)

	ctx, err := h.middleware.Require(ctx, constant.GET)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	report, err := h.uc.GetBestSeller(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var bestBooks []*toko.BestBook
	for _, r := range report {
		bestBooks = append(bestBooks, &toko.BestBook{
			Id:    int32(r.Id),
			Title: r.Title,
		})
	}

	return &toko.BestSellerReport{
		Buku: bestBooks,
	}, nil
}
func (h *ReportHandler) GetPriceReport(ctx context.Context, req *toko.EmptyReport) (*toko.PriceReport, error) {
	h.log.Info("get sales report in handler", req)

	ctx, err := h.middleware.Require(ctx, constant.GET)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	report, err := h.uc.GetPriceReport(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &toko.PriceReport{
		Max: report.Max,
		Min: report.Min,
		Avg: report.Avg,
	}, nil
}
