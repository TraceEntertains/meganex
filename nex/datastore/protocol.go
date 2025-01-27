package datastore

import (
	commondatastore "github.com/PretendoNetwork/nex-protocols-common-go/v2/datastore"
)

func NewDatastoreProtocol(protocol *commondatastore.CommonProtocol) error {
	err := initDatabase()
	if err != nil {
		return err
	}

	protocol.InitializeObjectByPreparePostParam = InitializeObjectByPreparePostParam
	protocol.InitializeObjectRatingWithSlot = InitializeObjectRatingWithSlot

	protocol.GetObjectInfoByPersistenceTargetWithPassword = GetObjectInfoByPersistenceTargetWithPassword
	protocol.GetObjectInfoByDataIDWithPassword = GetObjectInfoByDataIDWithPassword
	return nil
}

// func GetObjectInfoByDataID(dataID types.UInt64) (datastore_types.DataStoreMetaInfo, *nex.Error) {

// }
// func UpdateObjectPeriodByDataIDWithPassword(dataID types.UInt64, dataType types.UInt16, password types.UInt64) *nex.Error {

// }
// func UpdateObjectMetaBinaryByDataIDWithPassword(dataID types.UInt64, metaBinary types.QBuffer, password types.UInt64) *nex.Error {

// }
// func UpdateObjectDataTypeByDataIDWithPassword(dataID types.UInt64, period types.UInt16, password types.UInt64) *nex.Error {

// }
// func GetObjectSizeByDataID(dataID types.UInt64) (uint32, *nex.Error) {

// }
// func UpdateObjectUploadCompletedByDataID(dataID types.UInt64, uploadCompleted bool) *nex.Error {

// }
// func GetObjectInfoByDataIDWithPassword(dataID types.UInt64, password types.UInt64) (datastore_types.DataStoreMetaInfo, *nex.Error) {

// }
// func InitializeObjectByPreparePostParam(ownerPID types.PID, param datastore_types.DataStorePreparePostParam) (uint64, *nex.Error) {

// }
// func InitializeObjectRatingWithSlot(dataID uint64, param datastore_types.DataStoreRatingInitParamWithSlot) *nex.Error {

// }
// func RateObjectWithPassword(dataID types.UInt64, slot types.UInt8, ratingValue types.Int32, accessPassword types.UInt64) (datastore_types.DataStoreRatingInfo, *nex.Error) {

// }
// func DeleteObjectByDataIDWithPassword(dataID types.UInt64, password types.UInt64) *nex.Error {

// }
// func DeleteObjectByDataID(dataID types.UInt64) *nex.Error {

// }
// func GetObjectInfosByDataStoreSearchParam(param datastore_types.DataStoreSearchParam, pid types.PID) ([]datastore_types.DataStoreMetaInfo, uint32, *nex.Error) {

// }
// func GetObjectOwnerByDataID(dataID types.UInt64) (uint32, *nex.Error) {

// }
