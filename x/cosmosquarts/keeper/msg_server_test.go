package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/dmdv/cosmos-quarts/testutil/keeper"
	"github.com/dmdv/cosmos-quarts/x/cosmosquarts/keeper"
	"github.com/dmdv/cosmos-quarts/x/cosmosquarts/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.CosmosquartsKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
