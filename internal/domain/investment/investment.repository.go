package investment

import (
	"github.com/BagusAK95/amarta_test/internal/domain/common/repository"
)

type IInvestmentRepository interface {
	repository.IBaseRepo[Investment]
}
