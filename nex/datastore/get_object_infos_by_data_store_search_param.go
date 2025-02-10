package datastore

import (
	"database/sql"
	"errors"
	"meganex/globals"
	"slices"

	datastore_constants "meganex/nex/datastore/constants"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastoretypes "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/lib/pq"
)

var selectObjectsBySearchParamStmt *sql.Stmt

func GetObjectInfosByDataStoreSearchParam(param datastoretypes.DataStoreSearchParam, pid types.PID) ([]datastoretypes.DataStoreMetaInfo, uint32, *nex.Error) {
	if param.CreatedAfter == 0x9C3F3F7EFB { // 9999-12-31 23:59:59
		param.CreatedAfter = 0x4420000 // year 1, month 1, day 1 00:00:00 (go zero date)
	}

	if param.UpdatedAfter == 0x9C3F3F7EFB {
		param.UpdatedAfter = 0x4420000
	}

	if globals.NexConfig.DatastoreTrace {
		globals.Logger.Infof("param: %s\npid: %s", pid.String(), param.FormatToString(0))
	}

	tagArray := make([]string, len(param.Tags))
	for i := 0; i < len(param.Tags); i++ {
		tagArray[i] = param.Tags[i].String()
	}

	var idArray []int64

	if uint8(param.SearchTarget) == uint8(datastore_constants.SearchTypeFriend) {
		pids := globals.GetUserFriendPIDs(uint32(pid))

		// this is guessed behavior, it probably is just filtered to friends only with param.OwnerIDs ignored on official servers but no evidence for now
		// if we arent trying to filter then copy over the pids to idArray
		if len(param.OwnerIDs) == 0 {
			idArray = make([]int64, len(pids))
			for i := range pids {
				idArray = append(idArray, int64(pids[i]))
			}
		} else {
			// otherwise if we are then filter to pids present in param.OwnerIDs
			idArray = make([]int64, 0)
			for i := range pids {
				if slices.Contains(param.OwnerIDs, types.NewPID(uint64(pids[i]))) {
					idArray = append(idArray, int64(pids[i]))
				}
			}
		}
	} else {
		idArray = make([]int64, len(param.OwnerIDs))
		for i := range param.OwnerIDs {
			idArray = append(idArray, int64(param.OwnerIDs[i]))
		}
	}

	if uint8(param.SearchTarget) == uint8(datastore_constants.SearchTypeOwnAll) {
		idArray = append(idArray, int64(pid))
	}

	objects, err := getObjects(selectObjectsBySearchParamStmt,
		pq.Int64Array(idArray),
		uint16(param.DataType),
		uint8(param.SearchTarget),
		pq.FormatTimestamp(param.CreatedAfter.Standard()),
		pq.FormatTimestamp(param.CreatedBefore.Standard()),
		pq.FormatTimestamp(param.UpdatedAfter.Standard()),
		pq.FormatTimestamp(param.UpdatedBefore.Standard()),
		uint32(param.ReferDataID),
		pq.StringArray(tagArray),
		uint32(param.ResultRange.Length))

	if errors.Is(err, sql.ErrNoRows) {
		return []datastoretypes.DataStoreMetaInfo{}, 0, nil
	} else if err != nil {
		return []datastoretypes.DataStoreMetaInfo{}, 0, nex.NewError(nex.ResultCodes.DataStore.SystemFileError, err.Error())
	}

	if globals.NexConfig.DatastoreTrace {
		//globals.Logger.Infof("result: %v", objects)
		globals.Logger.Infof("object count: %v", len(objects))
	}

	return objects, uint32(len(objects)), nil
}

func initSelectObjectsBySearchParamStmt() error {
	stmt, err := Database.Prepare(
		selectObject + ` WHERE (owner = ANY($1) OR cardinality($1) = 0)
		AND deleted = 'false'
		AND data_type = $2
		AND CASE
		WHEN $3=1 THEN permission = 0
		WHEN $3=6 THEN permission = 1
		ELSE true END
		AND creation_date BETWEEN $4 AND $5
		AND update_date BETWEEN $6 AND $7
		AND refer_data_id = $8
		AND tags @> $9
		ORDER BY data_id DESC 
		LIMIT $10`)
	if err != nil {
		return err
	}

	selectObjectsBySearchParamStmt = stmt
	return nil
}
