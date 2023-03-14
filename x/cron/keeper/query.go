package keeper

import (
	"github.com/dmdv/cosmos-quarts/x/cron/types"
)

var _ types.QueryServer = Keeper{}
