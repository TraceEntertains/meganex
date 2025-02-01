package datastore

import (
	commondatastore "github.com/PretendoNetwork/nex-protocols-common-go/v2/datastore"
	"meganex/globals"
)

func NewDatastoreProtocol(protocol *commondatastore.CommonProtocol) error {
	err := initDatabase()
	if err != nil {
		return err
	}

	if globals.MinIOClient != nil {
		protocol.SetMinIOClient(globals.MinIOClient)
		protocol.S3Bucket = globals.NexConfig.DatastoreS3Bucket

		protocol.SetDataKeyBase("object")
		protocol.SetNotifyKeyBase("notif")
	}

	protocol.InitializeObjectByPreparePostParam = InitializeObjectByPreparePostParam
	protocol.InitializeObjectRatingWithSlot = InitializeObjectRatingWithSlot

	protocol.GetObjectInfoByDataID = GetObjectInfoByDataID
	protocol.GetObjectInfoByPersistenceTargetWithPassword = GetObjectInfoByPersistenceTargetWithPassword
	protocol.GetObjectInfoByDataIDWithPassword = GetObjectInfoByDataIDWithPassword

	protocol.GetObjectOwnerByDataID = GetObjectOwnerByDataID
	protocol.GetObjectSizeByDataID = GetObjectSizeByDataID

	protocol.UpdateObjectUploadCompletedByDataID = UpdateObjectUploadCompletedByDataID
	protocol.UpdateObjectMetaBinaryByDataIDWithPassword = UpdateObjectMetaBinaryByDataIDWithPassword
	protocol.UpdateObjectPeriodByDataIDWithPassword = UpdateObjectPeriodByDataIDWithPassword
	protocol.UpdateObjectDataTypeByDataIDWithPassword = UpdateObjectDataTypeByDataIDWithPassword

	protocol.DeleteObjectByDataID = DeleteObjectByDataID
	protocol.DeleteObjectByDataIDWithPassword = DeleteObjectByDataIDWithPassword
	return nil
}

// func RateObjectWithPassword(dataID types.UInt64, slot types.UInt8, ratingValue types.Int32, accessPassword types.UInt64) (datastore_types.DataStoreRatingInfo, *nex.Error) {

// }
// func GetObjectInfosByDataStoreSearchParam(param datastore_types.DataStoreSearchParam, pid types.PID) ([]datastore_types.DataStoreMetaInfo, uint32, *nex.Error) {

// }
