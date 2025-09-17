package repository

import (
	"context"

	"github.com/BagusAK95/amarta_test/internal/domain/common/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: Refactor using sq
type IBaseRepo[M model.EntityModel] interface {
	GetAll(ctx context.Context) ([]M, error)
	GetByID(ctx context.Context, ID uuid.UUID) (M, error)
	GetByIDs(ctx context.Context, IDs []uuid.UUID) ([]M, error)
	Create(ctx context.Context, model M) (M, error)
	CreateBulk(ctx context.Context, models []M) error
	Update(ctx context.Context, ID uuid.UUID, model M) (M, error)
	UpdateBulk(ctx context.Context, IDs []uuid.UUID, payload map[string]any) error
	UpdateWithMap(ctx context.Context, ID uuid.UUID, payload map[string]any) (M, error)
	Delete(ctx context.Context, ID uuid.UUID) error
	DeleteBulk(ctx context.Context, IDs []uuid.UUID) error
	CreateWithTx(ctx context.Context, model M, trx *gorm.DB) (M, error)
	CreateBulkWithTx(ctx context.Context, models []M, trx *gorm.DB) error
	CreateBulkAndReturnWithTx(ctx context.Context, models []M, trx *gorm.DB) ([]M, error)
	UpdateWithTx(ctx context.Context, ID uuid.UUID, model M, trx *gorm.DB) (M, error)
	UpdateBulkWithTx(ctx context.Context, IDs []uuid.UUID, payload map[string]any, trx *gorm.DB) error
	UpdateWithMapTx(ctx context.Context, ID uuid.UUID, payload map[string]any, trx *gorm.DB) (M, error)
	DeleteWithTx(ctx context.Context, ID uuid.UUID, trx *gorm.DB) error
	DeleteBulkWithTx(ctx context.Context, IDs []uuid.UUID, trx *gorm.DB) error
	BeginTransaction(ctx context.Context) *gorm.DB
	Rollback(trx *gorm.DB) *gorm.DB
	Commit(trx *gorm.DB) *gorm.DB
}

type BaseRepo[M model.EntityModel] struct {
	Entity    M
	writeConn *gorm.DB
	readConn  *gorm.DB
}

func NewBaseRepo[M model.EntityModel](dbMaster *gorm.DB, dbSlave *gorm.DB) *BaseRepo[M] {
	return &BaseRepo[M]{
		writeConn: dbMaster,
		readConn:  dbSlave,
	}
}

func (r *BaseRepo[M]) GetAll(ctx context.Context) ([]M, error) {
	var models []M
	err := r.readConn.WithContext(ctx).Where("deleted_at IS NULL").Find(&models).Error
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (r *BaseRepo[M]) GetByID(ctx context.Context, ID uuid.UUID) (M, error) {
	var model M

	builder := sq.Select("*").From(model.TableName()).Where(sq.Eq{"id": ID}).Where("deleted_at IS NULL")
	qry, args, err := builder.ToSql()
	if err != nil {
		return model, err
	}

	err = r.readConn.WithContext(ctx).Raw(qry, args...).Scan(&model).Error
	if err != nil {
		return model, err
	}

	return model, nil
}

func (r *BaseRepo[M]) GetByIDs(ctx context.Context, IDs []uuid.UUID) (models []M, err error) {
	var model M

	builder := sq.Select("*").From(model.TableName()).Where(sq.Eq{"id": IDs}).Where("deleted_at IS NULL")
	qry, args, err := builder.ToSql()
	if err != nil {
		return models, err
	}

	err = r.readConn.WithContext(ctx).Raw(qry, args...).Scan(&models).Error
	if err != nil {
		return models, err
	}

	return models, nil
}

// Create execute a single insert without specified transaction
func (r *BaseRepo[M]) Create(ctx context.Context, model M) (M, error) {
	err := r.writeConn.WithContext(ctx).Table(model.TableName()).Create(&model).Error
	if err != nil {
		return model, err
	}

	return model, nil
}

// CreateWithTx execute a single insert with specified transaction
func (r *BaseRepo[M]) CreateWithTx(ctx context.Context, model M, trx *gorm.DB) (M, error) {
	err := trx.Create(&model).Error
	if err != nil {
		return model, err
	}

	return model, nil
}

// CreateBulk execute a bulk insert without specified transaction
func (r *BaseRepo[M]) CreateBulk(ctx context.Context, models []M) error {
	err := r.writeConn.WithContext(ctx).CreateInBatches(models, InsertBatchSize).Error
	if err != nil {
		return err
	}

	return nil
}

// CreateBulkWithTx execute a bulk insert without specified transaction
func (r *BaseRepo[M]) CreateBulkWithTx(ctx context.Context, models []M, trx *gorm.DB) error {
	err := trx.WithContext(ctx).CreateInBatches(models, InsertBatchSize).Error
	if err != nil {
		return err
	}

	return nil
}

// CreateBulkAndReturnWithTx execute a bulk insert with specified transaction and return the models
func (r *BaseRepo[M]) CreateBulkAndReturnWithTx(ctx context.Context, models []M, trx *gorm.DB) ([]M, error) {
	err := trx.WithContext(ctx).CreateInBatches(models, InsertBatchSize).Error
	if err != nil {
		return models, err
	}

	return models, nil
}

// Update execute bulk update without specified transaction
func (r *BaseRepo[M]) Update(ctx context.Context, ID uuid.UUID, model M) (M, error) {
	err := r.writeConn.WithContext(ctx).Model(&model).Where("id=?", ID).Updates(model).Scan(&model).Error
	if err != nil {
		return model, err
	}

	return model, nil
}

// UpdateWithTx execute a single update with specified transaction
func (r *BaseRepo[M]) UpdateWithTx(ctx context.Context, ID uuid.UUID, model M, trx *gorm.DB) (M, error) {
	err := trx.WithContext(ctx).Model(&model).Where("id=?", ID).Updates(model).Scan(&model).Error
	if err != nil {
		return model, err
	}

	return model, nil
}

// UpdateBulk execute a bulk update without specified transaction
func (r *BaseRepo[M]) UpdateBulk(ctx context.Context, IDs []uuid.UUID, payload map[string]any) error {
	var model M

	err := r.writeConn.WithContext(ctx).Model(&model).Where("id IN ?", IDs).Updates(payload).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateBulkWithTx execute a bulk update without specified transaction
func (r *BaseRepo[M]) UpdateBulkWithTx(ctx context.Context, IDs []uuid.UUID, payload map[string]any, trx *gorm.DB) error {
	var model M

	err := trx.WithContext(ctx).Model(&model).Where("id IN ?", IDs).Updates(payload).Error
	if err != nil {
		return err
	}
	return nil

}

// UpdateWithMap execute a single update with Map without specified transaction
func (r *BaseRepo[M]) UpdateWithMap(ctx context.Context, ID uuid.UUID, payload map[string]any) (M, error) {
	var model M

	err := r.writeConn.WithContext(ctx).Model(&model).Where("id=?", ID).Updates(payload).Scan(&model).Error
	if err != nil {
		return model, err
	}

	return model, nil
}

// UpdateWithMapTx execute a single update with Map with specified transaction
func (r *BaseRepo[M]) UpdateWithMapTx(ctx context.Context, ID uuid.UUID, payload map[string]any, trx *gorm.DB) (M, error) {
	var model M

	err := trx.WithContext(ctx).Model(&model).Where("id=?", ID).Updates(payload).Scan(&model).Error
	if err != nil {
		return model, err
	}

	return model, nil
}

// Delete execute a single delete without specified transaction
func (r *BaseRepo[M]) Delete(ctx context.Context, ID uuid.UUID) error {
	model, err := r.GetByID(ctx, ID)
	if err != nil {
		return err
	}

	err = r.writeConn.WithContext(ctx).Delete(&model, ID).Error
	if err != nil {
		return err
	}

	return nil
}

// Delete execute a single delete with specified transaction
func (r *BaseRepo[M]) DeleteWithTx(ctx context.Context, ID uuid.UUID, trx *gorm.DB) error {
	model, err := r.GetByID(ctx, ID)
	if err != nil {
		return err
	}

	err = trx.WithContext(ctx).Delete(&model).Error
	if err != nil {
		return err
	}

	return nil
}

// DeleteBulk execute a `soft` bulk delete without specified transaction
func (r *BaseRepo[M]) DeleteBulk(ctx context.Context, IDs []uuid.UUID) error {
	models, err := r.GetByIDs(ctx, IDs)
	if err != nil {
		return err
	}

	err = r.writeConn.WithContext(ctx).Delete(&models).Error
	if err != nil {
		return err
	}

	return nil
}

// DeleteBulk execute a bulk `soft` delete with specified transaction
func (r *BaseRepo[M]) DeleteBulkWithTx(ctx context.Context, IDs []uuid.UUID, trx *gorm.DB) error {
	models, err := r.GetByIDs(ctx, IDs)
	if err != nil {
		return err
	}

	err = trx.WithContext(ctx).Delete(&models).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepo[M]) BeginTransaction(ctx context.Context) *gorm.DB {
	return r.writeConn.WithContext(ctx).Begin()
}
func (r *BaseRepo[M]) Rollback(trx *gorm.DB) *gorm.DB {
	return trx.Rollback()
}
func (r *BaseRepo[M]) Commit(trx *gorm.DB) *gorm.DB {
	return trx.Commit()
}
