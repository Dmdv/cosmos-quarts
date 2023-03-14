package keeper

import (
	"errors"
	"fmt"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	it "github.com/dmdv/cosmos-quarts/x/quartz/types"
)

// Keeper struct
type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	scheduleChan chan it.Schedule
}

// NewKeeper creates a new instance of the cron Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, scheduleChan chan it.Schedule) Keeper {
	return Keeper{
		storeKey:     storeKey,
		cdc:          cdc,
		scheduleChan: scheduleChan,
	}
}

// AddSchedule adds a new scheduled task to the cron module
func (k Keeper) AddSchedule(ctx sdk.Context, schedule it.Schedule) error {
	store := ctx.KVStore(k.storeKey)

	// Encode and store the schedule
	bz := k.cdc.MustMarshalBinaryBare(schedule)
	store.Set([]byte(schedule.Schedule), bz)

	// Notify the scheduler of the new schedule
	k.scheduleChan <- schedule

	return nil
}

// GetSchedule returns the schedule with the given ID
func (k Keeper) GetSchedule(ctx sdk.Context, scheduleID string) (it.Schedule, error) {
	store := ctx.KVStore(k.storeKey)

	// Get the encoded schedule from the store
	bz := store.Get([]byte(scheduleID))
	if bz == nil {
		return it.Schedule{}, errors.New("schedule not found")
	}

	// Decode the schedule and return
	var schedule it.Schedule
	k.cdc.MustUnmarshalBinaryBare(bz, &schedule)
	return schedule, nil
}

// GetAllSchedules returns all schedules in the store
func (k Keeper) GetAllSchedules(ctx sdk.Context) []it.Schedule {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)

	// Iterate over all schedules and decode them
	var schedules []it.Schedule
	for ; iterator.Valid(); iterator.Next() {
		var schedule it.Schedule
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &schedule)
		schedules = append(schedules, schedule)
	}

	return schedules
}

// DeleteSchedule deletes the schedule with the given ID
func (k Keeper) DeleteSchedule(ctx sdk.Context, scheduleID string) error {
	store := ctx.KVStore(k.storeKey)

	// Delete the schedule from the store
	store.Delete([]byte(scheduleID))

	return nil
}

// ProcessSchedules checks all schedules and executes any that are due
func (k Keeper) ProcessSchedules(ctx sdk.Context) {
	// Get the current block time
	blockTime := ctx.BlockTime()

	// Iterate over all schedules and execute any that are due
	for _, schedule := range k.GetAllSchedules(ctx) {
		nextTime, err := schedule.Next(blockTime)
		if err != nil {
			// Log error and continue to next schedule
			fmt.Printf("Error processing schedule %s: %v\\n", schedule.Schedule, err)
			continue
		}

		if nextTime.After(blockTime) {
			// Schedule is not due yet, skip it
			continue
		}

		// Schedule is due, execute it
		for _, msg := range schedule.Msgs {
			if err := sdk.ValidateBasic(msg); err != nil {
				// Log error and continue to next message
				fmt.Printf("Error processing message: %v\\n", err)
				continue
			}

			_, err := sdk.ApplyAndReturnTxs(ctx, []sdk.Msg{msg})
			if err != nil {
				// Log error and continue to next message
				fmt.Printf("Error processing message: %v\\n", err)
				continue
			}
		}

		// Delete the schedule now that it has been executed
		k.DeleteSchedule(ctx, schedule.Schedule)
	}
}

// HandleMsgSchedule handles a MsgSchedule message
func HandleMsgSchedule(ctx sdk.Context, keeper Keeper, msg MsgSchedule) (*sdk.Result, error) {
	schedule, err := it.NewSchedule(msg.Schedule, ctx.BlockHeight(), msg.Msgs, msg.Labels)
	if err != nil {
		return nil, err
	}

	if err := keeper.AddSchedule(ctx, schedule); err != nil {
		return nil, err
	}

	return &sdk.Result{
		Events: sdk.Events{
			sdk.NewEvent(
				EventTypeSchedule,
				sdk.NewAttribute(AttributeKeyScheduleID, schedule.Schedule),
			),
		},
	}, nil
}

// HandleMsgUnschedule handles a MsgUnschedule message
func HandleMsgUnschedule(ctx sdk.Context, keeper Keeper, msg MsgUnschedule) (*sdk.Result, error) {
	if err := keeper.DeleteSchedule(ctx, msg.ScheduleID); err != nil {
		return nil, err
	}

	return &sdk.Result{
		Events: sdk.Events{
			sdk.NewEvent(
				EventTypeUnschedule,
				sdk.NewAttribute(AttributeKeyScheduleID, msg.ScheduleID),
			),
		},
	}, nil
}

// Querier struct
type Querier struct {
	keeper Keeper
}

// NewQuerier creates a new instance of the cron Querier
func NewQuerier(keeper Keeper) Querier {
	return Querier{keeper: keeper}
}

// QuerySchedule returns the schedule with the given ID
func (q Querier) QuerySchedule(ctx sdk.Context, req abci.RequestQuery) ([]byte, error) {
	var params QueryScheduleParams
	if err := q.keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, err
	}

	schedule, err := q.keeper.GetSchedule(ctx, params.ScheduleID)
	if err != nil {
		return nil, err
	}

	bz, err := q.keeper.cdc.MarshalJSON(schedule)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

// QueryAllSchedules returns all schedules in the store
func (q Querier) QueryAllSchedules(ctx sdk.Context, req abci.RequestQuery) ([]byte, error) {
	schedules := q.keeper.GetAllSchedules(ctx)

	bz, err := q.keeper.cdc.MarshalJSON(schedules)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

// QueryParams returns the current module parameters
func (q Querier) QueryParams(ctx sdk.Context, req abci.RequestQuery) ([]byte, error) {
	return q.keeper.cdc.MarshalJSON(q.keeper.GetParams(ctx))
}
