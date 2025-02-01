package datastore

import (
	"database/sql"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastoretypes "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/lib/pq"
	"meganex/globals"
	"time"
)

var Database *sql.DB

func initDatabase() error {
	inits := []func() error{
		initTables,
		initInsertObjectStmt,                   // initialize_object_by_prepare_post_param.go
		initSelectObjectByIdPasswordStmt,       // get_object_info_by_data_id_with_password.go
		initSelectObjectByOwnerPersistenceStmt, // get_object_info_by_persistence_target_with_password.go
		initSelectObjectByIdStmt,               // get_object_info_by_data_id.go
		initSelectOwnerByIdStmt,                // get_object_owner_by_data_id.go
		initSelectSizeByIdStmt,                 // get_object_size_by_data_id.go
		initUpdateUploadCompleteByIdStmt,       // update_object_upload_completed_by_data_id.go
		initUpdateMetaBinaryByIdPasswordStmt,   // update_object_meta_binary_by_data_id_with_password.go
		initUpdatePeriodByIdPasswordStmt,       // update_object_period_by_data_id_with_password.go
		initUpdateDataTypeByIdPasswordStmt,     // update_object_data_type_by_data_id_with_password.go
		initUpdateDeletedByIdStmt,              // delete_object_by_data_id.go
		initUpdateDeletedByIdPasswordStmt,      // delete_object_by_data_id_with_password.go
	}

	for _, init := range inits {
		err := init()
		if err != nil {
			return err
		}
	}

	return nil
}

func initTables() error {
	_, err := Database.Exec(`CREATE SCHEMA IF NOT EXISTS datastore`)
	if err != nil {
		return err
	}

	globals.Logger.Success("datastore Postgres schema created")

	_, err = Database.Exec(`CREATE SEQUENCE IF NOT EXISTS datastore.object_data_id_seq
		INCREMENT 1
		MINVALUE 1
		MAXVALUE 281474976710656
		START 1
		CACHE 1`, // * Honestly I don't know what CACHE does but I saw it recommended so here it is
	)
	if err != nil {
		return err
	}

	_, err = Database.Exec(`CREATE TABLE IF NOT EXISTS datastore.objects (
		data_id bigint NOT NULL DEFAULT nextval('datastore.object_data_id_seq') PRIMARY KEY,
		upload_completed boolean NOT NULL DEFAULT FALSE,
		deleted boolean NOT NULL DEFAULT FALSE,
		owner bigint,
		size int,
		name text,
		data_type int,
		meta_binary bytea,
		permission int,
		permission_recipients int[],
		delete_permission int,
		delete_permission_recipients int[],
		flag int,
		period int,
		refer_data_id bigint,
		tags text[],
		persistence_slot_id int,
		extra_data text[],
		access_password bigint NOT NULL DEFAULT 0,
		update_password bigint NOT NULL DEFAULT 0,
		creation_date timestamp,
		update_date timestamp
	)`)
	if err != nil {
		return err
	}

	globals.Logger.Success("Postgres tables created")
	return nil
}

const selectObject = `
	SELECT
		data_id,
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
		creation_date,
		update_date
	FROM datastore.objects`

// Helper to unpack things selected with (selectObject + ` WHERE ....`)
func getObjects(stmt *sql.Stmt, args ...any) ([]datastoretypes.DataStoreMetaInfo, error) {
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	// surely we know the length of the result set at this point?
	var results []datastoretypes.DataStoreMetaInfo

	for rows.Next() {
		result := datastoretypes.NewDataStoreMetaInfo()

		var createdTime time.Time
		var updatedTime time.Time

		// TODO check this - it's stolen from SMM DataStore but seems fifty shades of fucked up for a generic impl
		result.ExpireTime = types.NewDateTime(0x9C3f3E0000) // * 9999-12-31T00:00:00.000Z. This is what the real server sends

		err := rows.Scan(
			&result.DataID,
			&result.OwnerID,
			&result.Size,
			&result.Name,
			&result.DataType,
			&result.MetaBinary,
			&result.Permission.Permission,
			pq.Array(&result.Permission.RecipientIDs),
			&result.DelPermission.Permission,
			pq.Array(&result.DelPermission.RecipientIDs),
			&result.Flag,
			&result.Period,
			&result.ReferDataID,
			pq.Array(&result.Tags),
			&createdTime,
			&updatedTime,
		)
		if err != nil {
			return nil, err
			//globals.Logger.Error(err.Error())
			//continue
		}

		// I'm not sure how this API is meant to be used but this works
		result.CreatedTime = result.CreatedTime.FromTimestamp(createdTime)
		result.UpdatedTime = result.UpdatedTime.FromTimestamp(updatedTime)
		result.ReferredTime = result.ReferredTime.FromTimestamp(createdTime)

		results = append(results, result)
	}

	return results, rows.Err()
}
