package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dmdv/cosmos-quarts/x/cron/types"
)

func (k Keeper) ScheduleTask(
	ctx sdk.Context,
	spec string,
	taskFunc types.Function,
	settings map[string]interface{},
) error {
	return k.cronScheduler.RegisterTask(ctx, spec, taskFunc, settings)
}
