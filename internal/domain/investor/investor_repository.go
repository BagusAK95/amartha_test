package investor

import (
	"github.com/BagusAK95/amarta_test/internal/domain/common/repository"
)

type IInvestorRepository interface {
	repository.IBaseRepo[Investor]
}
