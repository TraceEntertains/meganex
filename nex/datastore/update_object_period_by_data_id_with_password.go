package datastore

import (
	"database/sql"
	"errors"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"meganex/globals"
	"time"
)

var updatePeriodByIdPasswordStmt *sql.Stmt

func UpdateObjectPeriodByDataIDWithPassword(dataID types.UInt64, period types.UInt16, password types.UInt64) *nex.Error {
	if globals.NexConfig.DatastoreTrace {
		globals.Logger.Infof("dataID: %v\nperiod: %v\npassword: %v", dataID, period, password)
	}

	result, err := updatePeriodByIdPasswordStmt.Exec(dataID, password, period, time.Now())
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

func initUpdatePeriodByIdPasswordStmt() error {
	stmt, err := Database.Prepare(`UPDATE datastore.objects SET period = $3, update_date = $4 WHERE data_id = $1 AND update_password = $2`)
	if err != nil {
		return err
	}

	updatePeriodByIdPasswordStmt = stmt
	return nil
}
