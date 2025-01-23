package nex

import (
	"meganex/globals"

	"github.com/PretendoNetwork/nex-go/v2/types"

	commonsecure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"

	commonnattraversal "github.com/PretendoNetwork/nex-protocols-common-go/v2/nat-traversal"
	nattraversal "github.com/PretendoNetwork/nex-protocols-go/v2/nat-traversal"

	commonmatchmaking "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making"
	matchmaking "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"

	commonmatchmakingext "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making-ext"
	matchmakingext "github.com/PretendoNetwork/nex-protocols-go/v2/match-making-ext"

	commonmatchmakeextension "github.com/PretendoNetwork/nex-protocols-common-go/v2/matchmake-extension"
	matchmakeextension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"

	megadatastore "meganex/nex/datastore"

	commondatastore "github.com/PretendoNetwork/nex-protocols-common-go/v2/datastore"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
)

func registerCommonSecureServerProtocols() {
	var failed = 0
	for _, protocol := range globals.NexConfig.SecureProtocols {
		switch protocol {

		case "secure":
			secureProtocol := secure.NewProtocol()
			globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
			commonSecureProtocol := commonsecure.NewCommonProtocol(secureProtocol)

			// Stubbed until we can parse these reports
			commonSecureProtocol.CreateReportDBRecord = func(_ types.PID, _ types.UInt32, _ types.QBuffer) error {
				return nil
			}

		case "natTraversal":
			natTraversalProtocol := nattraversal.NewProtocol()
			globals.SecureEndpoint.RegisterServiceProtocol(natTraversalProtocol)
			commonnattraversal.NewCommonProtocol(natTraversalProtocol)

		case "matchMaking":
			matchMakingProtocol := matchmaking.NewProtocol()
			globals.SecureEndpoint.RegisterServiceProtocol(matchMakingProtocol)
			commonMatchMakingProtocol := commonmatchmaking.NewCommonProtocol(matchMakingProtocol)
			commonMatchMakingProtocol.SetManager(globals.MatchmakingManager)

		case "matchMakingExt":
			matchMakingExtProtocol := matchmakingext.NewProtocol()
			globals.SecureEndpoint.RegisterServiceProtocol(matchMakingExtProtocol)
			commonMatchMakingExtProtocol := commonmatchmakingext.NewCommonProtocol(matchMakingExtProtocol)
			commonMatchMakingExtProtocol.SetManager(globals.MatchmakingManager)

		case "matchmakeExtension":
			matchmakeExtensionProtocol := matchmakeextension.NewProtocol()
			globals.SecureEndpoint.RegisterServiceProtocol(matchmakeExtensionProtocol)
			commonMatchmakeExtensionProtocol := commonmatchmakeextension.NewCommonProtocol(matchmakeExtensionProtocol)
			commonMatchmakeExtensionProtocol.SetManager(globals.MatchmakingManager)

		case "datastore":
			datastoreProtocol := datastore.NewProtocol()
			globals.SecureEndpoint.RegisterServiceProtocol(datastoreProtocol)
			commonDatastoreProtocol := commondatastore.NewCommonProtocol(datastoreProtocol)
			megadatastore.Database = globals.Postgres
			megadatastore.NewDatastoreProtocol(commonDatastoreProtocol)

		default:
			globals.Logger.Warningf("Ignoring unknown protocol \"%v\"!", protocol)
			failed++
		}
	}

	globals.Logger.Infof("Configured %v protocols.", len(globals.NexConfig.SecureProtocols)-failed)
}
