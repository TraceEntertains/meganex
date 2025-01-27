package datastore

import (
	"database/sql"
	"meganex/globals"
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	types "github.com/PretendoNetwork/nex-go/v2/types"
	datastoretypes "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/lib/pq"
)

var insertObjectStmt *sql.Stmt

func InitializeObjectByPreparePostParam(ownerPID types.PID, param datastoretypes.DataStorePreparePostParam) (uint64, *nex.Error) {
	now := time.Now()

	globals.Logger.Infof("posting %v", param.FormatToString(0))

	var dataID uint64

	err := insertObjectStmt.QueryRow(
		ownerPID,
		param.Size,
		param.Name,
		param.DataType,
		param.MetaBinary,
		param.Permission.Permission,
		pq.Array(param.Permission.RecipientIDs),
		param.DelPermission.Permission,
		pq.Array(param.DelPermission.RecipientIDs),
		param.Flag,
		param.Period,
		param.ReferDataID,
		pq.Array(param.Tags),
		param.PersistenceInitParam.PersistenceSlotID, // todo DeleteLastObject
		pq.Array(param.ExtraData),
		now,
		now,
	).Scan(&dataID)
	if err != nil {
		return 0, nex.NewError(nex.ResultCodes.DataStore.SystemFileError, err.Error())
	}

	return dataID, nil
}

func initInsertObjectStmt() error {
	stmt, err := Database.Prepare(`INSERT INTO datastore.objects 
	(
		owner,
		size,
		name,
		data_type,
		meta_binary,
		permission,
		permission_recipients,
		delete_permission,
		delete_permission_recipients,
		flag,
		period,
		refer_data_id,
		tags,
		persistence_slot_id,
		extra_data,
		creation_date,
		update_date
	)
	VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
	) RETURNING data_id`)
	if err != nil {
		return err
	}

	insertObjectStmt = stmt
	return nil
}
