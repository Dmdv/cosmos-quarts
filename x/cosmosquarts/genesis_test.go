package cosmosquarts_test

import (
	"testing"

	keepertest "github.com/dmdv/cosmos-quarts/testutil/keeper"
	"github.com/dmdv/cosmos-quarts/testutil/nullify"
	"github.com/dmdv/cosmos-quarts/x/cosmosquarts"
	"github.com/dmdv/cosmos-quarts/x/cosmosquarts/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.CosmosquartsKeeper(t)
	cosmosquarts.InitGenesis(ctx, *k, genesisState)
	got := cosmosquarts.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
