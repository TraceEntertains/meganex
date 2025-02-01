package datastore

import (
	"database/sql"
	"errors"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"meganex/globals"
	"time"
)

var updateDeletedByIdStmt *sql.Stmt

func DeleteObjectByDataID(dataID types.UInt64) *nex.Error {
	if globals.NexConfig.DatastoreTrace {
		globals.Logger.Infof("dataID: %v", dataID)
	}

	result, err := updateDeletedByIdStmt.Exec(dataID, true, time.Now())
	if errors.Is(err, sql.ErrNoRows) {
		return nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found or wrong password")
	} else if err != nil {
		return nex.NewError(nex.ResultCodes.DataStore.SystemFileError, err.Error())
	}

	rows, err := result.RowsAffected()
	if err != nil && rows < 1 {
		return nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found or wrong password")
	}

	return nil
}

func initUpdateDeletedByIdStmt() error {
	stmt, err := Database.Prepare(`UPDATE datastore.objects SET deleted = $2, update_date = $3 WHERE data_id = $1`)
	if err != nil {
		return err
	}

	updateDeletedByIdStmt = stmt
	return nil
}
