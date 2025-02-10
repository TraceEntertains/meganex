package datastore

import (
	"database/sql"
	"errors"
	"meganex/globals"

	"github.com/PretendoNetwork/nex-go/v2"
	datastoretypes "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
)

var selectObjectByOwnerPersistenceNoPassStmt *sql.Stmt

func GetObjectInfoByPersistenceTarget(persistenceTarget datastoretypes.DataStorePersistenceTarget) (datastoretypes.DataStoreMetaInfo, *nex.Error) {
	if globals.NexConfig.DatastoreTrace {
		globals.Logger.Infof("persistenceTarget: %v", persistenceTarget)
	}

	objects, err := getObjects(selectObjectByOwnerPersistenceNoPassStmt, persistenceTarget.OwnerID, persistenceTarget.PersistenceSlotID)
	if errors.Is(err, sql.ErrNoRows) || len(objects) < 1 {
		// todo nex.ResultCodes.DataStore.InvalidPassword return?
		return datastoretypes.NewDataStoreMetaInfo(), nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found or wrong password")
	} else if err != nil {
		return datastoretypes.NewDataStoreMetaInfo(), nex.NewError(nex.ResultCodes.DataStore.SystemFileError, err.Error())
	}

	if globals.NexConfig.DatastoreTrace {
		globals.Logger.Infof("result: %v", objects[0])
	}
	return objects[0], nil
}

func initSelectObjectByOwnerPersistenceNoPassStmt() error {
	stmt, err := Database.Prepare(selectObject + `
		WHERE owner = $1 AND persistence_slot_id = $2 AND deleted = 'false'
		LIMIT 1
	`)
	if err != nil {
		return err
	}

	selectObjectByOwnerPersistenceNoPassStmt = stmt
	return nil
}
