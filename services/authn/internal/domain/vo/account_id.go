package vo

import "github.com/bcchicr/cinemaroom/services/authn/internal/domain/common"

type accountIDTag struct{}

type AccountID = common.TypedId[accountIDTag]

func NewAccountIDFromString(rawId string) (*AccountID, error) {
	return common.NewTypeIdFromString[accountIDTag](rawId)
}
