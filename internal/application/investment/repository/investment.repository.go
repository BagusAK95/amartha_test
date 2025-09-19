package repository

import (
	"context"

	"github.com/BagusAK95/amarta_test/internal/domain/common/repository"
	"github.com/BagusAK95/amarta_test/internal/domain/investment"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

var tracerName = "InvestmentRepository"
var tracer = otel.Tracer(tracerName)

type investmentRepo struct {
	repository.BaseRepo[investment.Investment]
	writeConn *gorm.DB
	readConn  *gorm.DB
}

func NewInvestmentRepo(dbMaster *gorm.DB, dbSlave *gorm.DB) investment.IInvestmentRepository {
	baseRepo := repository.NewBaseRepo[investment.Investment](dbMaster, dbSlave)

	return &investmentRepo{
		BaseRepo:  *baseRepo,
		writeConn: dbMaster,
		readConn:  dbSlave,
	}
}

func (r *investmentRepo) GetTotalInvestmentByLoanID(ctx context.Context, loanID uuid.UUID) (total float64, err error) {
	ctx, span := tracer.Start(ctx, tracerName+".GetTotalInvestmentByLoanID")
	defer span.End()

	var model investment.Investment

	builder := sq.
		Select("COALESCE(SUM(amount), 0)").
		From(model.TableName()).
		Where(sq.Eq{
			"loan_id":    loanID,
			"deleted_at": nil,
		})

	qry, args, err := builder.ToSql()
	if err != nil {
		return
	}

	err = r.readConn.WithContext(ctx).Raw(qry, args...).Scan(&total).Error
	if err != nil {
		return
	}

	return
}
