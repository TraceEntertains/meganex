package datastore

import (
	"github.com/PretendoNetwork/nex-go/v2"
	datastoretypes "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"meganex/globals"
)

func InitializeObjectRatingWithSlot(dataID uint64, param datastoretypes.DataStoreRatingInitParamWithSlot) *nex.Error {
	if globals.NexConfig.DatastoreTrace {
		globals.Logger.Infof("dataID: %v\nparam: %v", dataID, param)
	}

	return nex.NewError(nex.ResultCodes.Core.NotImplemented, "Ratings not yet implemented")
}
