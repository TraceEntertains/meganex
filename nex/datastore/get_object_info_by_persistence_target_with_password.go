package datastore

import (
	"database/sql"
	"errors"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastoretypes "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"meganex/globals"
)

var selectObjectByOwnerPersistenceStmt *sql.Stmt

func GetObjectInfoByPersistenceTargetWithPassword(persistenceTarget datastoretypes.DataStorePersistenceTarget, password types.UInt64) (datastoretypes.DataStoreMetaInfo, *nex.Error) {
	globals.Logger.Infof("target %v password %v", persistenceTarget.FormatToString(0), password)

	objects, err := getObjects(selectObjectByOwnerPersistenceStmt, persistenceTarget.OwnerID, persistenceTarget.PersistenceSlotID, password)
	if errors.Is(err, sql.ErrNoRows) || len(objects) < 1 {
		// todo nex.ResultCodes.DataStore.InvalidPassword return?
		return datastoretypes.NewDataStoreMetaInfo(), nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found or wrong password")
	} else if err != nil {
		return datastoretypes.NewDataStoreMetaInfo(), nex.NewError(nex.ResultCodes.DataStore.SystemFileError, err.Error())
	}

	globals.Logger.Infof("returning %v", objects[0])

	return objects[0], nil
}

func initSelectObjectByOwnerPersistenceStmt() error {
	// hack: order by data_id gets around the whole DeleteLastObject thing
	stmt, err := Database.Prepare(selectObject + `
		WHERE owner = $1 AND persistence_slot_id = $2 AND access_password = $3
		ORDER BY data_id DESC LIMIT 1
	`)
	if err != nil {
		return err
	}

	selectObjectByOwnerPersistenceStmt = stmt
	return nil
}
