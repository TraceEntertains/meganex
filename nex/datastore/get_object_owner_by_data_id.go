package datastore

import (
	"database/sql"
	"errors"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"meganex/globals"
)

var selectOwnerByIdStmt *sql.Stmt

func GetObjectOwnerByDataID(dataID types.UInt64) (uint32, *nex.Error) {
	if globals.NexConfig.DatastoreTrace {
		globals.Logger.Infof("dataID: %v", dataID)
	}

	var result uint32
	err := selectOwnerByIdStmt.QueryRow(dataID).Scan(&result)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found or wrong password")
	} else if err != nil {
		return 0, nex.NewError(nex.ResultCodes.DataStore.SystemFileError, err.Error())
	}

	if globals.NexConfig.DatastoreTrace {
		globals.Logger.Infof("result: %v", result)
	}
	return result, nil
}

func initSelectOwnerByIdStmt() error {
	stmt, err := Database.Prepare(`SELECT owner FROM datastore.objects WHERE data_id = $1`)
	if err != nil {
		return err
	}

	selectOwnerByIdStmt = stmt
	return nil
}
