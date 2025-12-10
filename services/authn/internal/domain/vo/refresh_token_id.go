package vo

import "github.com/bcchicr/cinemaroom/services/authn/internal/domain/common"

type refreshTokenIDTag struct{}

type RefreshTokenID = common.TypedId[refreshTokenIDTag]

func NewRefreshTokenIdFromString(rawId string) (*RefreshTokenID, error) {
	return common.NewTypeIdFromString[refreshTokenIDTag](rawId)
}
