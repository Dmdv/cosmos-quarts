package keeper_test

import (
	"testing"

	testkeeper "github.com/dmdv/cosmos-quarts/testutil/keeper"
	"github.com/dmdv/cosmos-quarts/x/cron/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.CronKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
