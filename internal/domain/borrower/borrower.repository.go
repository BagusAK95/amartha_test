package borrower

import (
	"github.com/BagusAK95/amarta_test/internal/domain/common/repository"
)

type IBorrowerRepository interface {
	repository.IBaseRepo[Borrower]
}
