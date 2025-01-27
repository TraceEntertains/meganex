package nex

import (
	"github.com/PretendoNetwork/nex-go/v2/types"
	commonmatchmakingext "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making-ext"
	matchmakingext "github.com/PretendoNetwork/nex-protocols-go/v2/match-making-ext"
	"meganex/globals"

	commonsecure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"

	commonnattraversal "github.com/PretendoNetwork/nex-protocols-common-go/v2/nat-traversal"
	nattraversal "github.com/PretendoNetwork/nex-protocols-go/v2/nat-traversal"

	commonmatchmaking "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making"
	matchmaking "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"

	commonmatchmakeextension "github.com/PretendoNetwork/nex-protocols-common-go/v2/matchmake-extension"
	matchmakeextension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"

	megadatastore "meganex/nex/datastore"

	commondatastore "github.com/PretendoNetwork/nex-protocols-common-go/v2/datastore"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
)

type ProtocolHandler struct {
	id   uint16
	init func()
}

var SecureProtocols = map[string]ProtocolHandler{
	"remoteLog": {1, nil},
	"natTraversal": {3, func() {
		natTraversalProtocol := nattraversal.NewProtocol()
		globals.SecureEndpoint.RegisterServiceProtocol(natTraversalProtocol)
		commonnattraversal.NewCommonProtocol(natTraversalProtocol)
	}},
	"ticketGranting": {10, nil},
	"secure": {11, func() {
		secureProtocol := secure.NewProtocol()
		globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
		commonSecureProtocol := commonsecure.NewCommonProtocol(secureProtocol)

		// Stubbed until we can parse these reports
		commonSecureProtocol.CreateReportDBRecord = func(_ types.PID, _ types.UInt32, _ types.QBuffer) error {
			return nil
		}
	}},
	"notifications": {14, nil},
	"health":        {18, nil},
	"monitoring":    {19, nil},
	"friends":       {20, nil},
	"matchMaking": {21, func() {
		matchMakingProtocol := matchmaking.NewProtocol()
		globals.SecureEndpoint.RegisterServiceProtocol(matchMakingProtocol)
		commonMatchMakingProtocol := commonmatchmaking.NewCommonProtocol(matchMakingProtocol)
		commonMatchMakingProtocol.SetManager(globals.MatchmakingManager)
	}},
	"messaging":         {23, nil},
	"persistentStore":   {24, nil},
	"accountManagement": {25, nil},
	"messageDelivery":   {27, nil},
	"matchMakingExt": {50, func() {
		matchMakingExtProtocol := matchmakingext.NewProtocol()
		globals.SecureEndpoint.RegisterServiceProtocol(matchMakingExtProtocol)
		commonMatchMakingExtProtocol := commonmatchmakingext.NewCommonProtocol(matchMakingExtProtocol)
		commonMatchMakingExtProtocol.SetManager(globals.MatchmakingManager)
	}},
	"nintendoNotifications": {100, nil},
	"friends3DS":            {101, nil},
	"friendsWiiU":           {102, nil},
	"matchmakeExtension": {109, func() {
		matchmakeExtensionProtocol := matchmakeextension.NewProtocol()
		globals.SecureEndpoint.RegisterServiceProtocol(matchmakeExtensionProtocol)
		commonMatchmakeExtensionProtocol := commonmatchmakeextension.NewCommonProtocol(matchmakeExtensionProtocol)
		commonMatchmakeExtensionProtocol.SetManager(globals.MatchmakingManager)
		commonMatchmakeExtensionProtocol.CleanupSearchMatchmakeSession = CleanupSearchMatchmakeSession
		commonMatchmakeExtensionProtocol.CleanupMatchmakeSessionSearchCriterias = CleanupMatchmakeSessionSearchCriterias
	}},
	"utility": {110, nil}, // todo: storagemanager?
	"ranking": {112, nil},
	"datastore": {115, func() {
		datastoreProtocol := datastore.NewProtocol()
		globals.SecureEndpoint.RegisterServiceProtocol(datastoreProtocol)
		commonDatastoreProtocol := commondatastore.NewCommonProtocol(datastoreProtocol)
		megadatastore.Database = globals.Postgres
		_ = megadatastore.NewDatastoreProtocol(commonDatastoreProtocol)
	}},
	"debug":            {116, nil},
	"subscription":     {117, nil},
	"serviceItem":      {119, nil},
	"matchmakeReferee": {120, nil},
	"subscriber":       {121, nil},
	"ranking2":         {122, nil},
	"aaUser":           {123, nil},
	"screening":        {124, nil},
}
var StartedSecureProtocols []uint16

func registerCommonSecureServerProtocols() {
	for _, name := range globals.NexConfig.SecureProtocols {
		protocol, ok := SecureProtocols[name]
		if !ok || protocol.init == nil {
			globals.Logger.Warningf("Skipping unknown/unimplemented protocol \"%v\"", name)
			continue
		}

		protocol.init()
		StartedSecureProtocols = append(StartedSecureProtocols, protocol.id)
	}

	started := make([]string, 0, len(StartedSecureProtocols))
	for _, id := range StartedSecureProtocols {
		name, _ := FindProtocolByID(id)
		started = append(started, name)
	}

	globals.Logger.Infof("Configured %v protocols: %v", len(StartedSecureProtocols), started)
}

func FindProtocolByID(id uint16) (string, *ProtocolHandler) {
	for k, v := range SecureProtocols {
		if v.id == id {
			return k, &v
		}
	}

	return "", nil
}
