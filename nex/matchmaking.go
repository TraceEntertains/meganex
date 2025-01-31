package nex

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	matchmakingtypes "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
	matchmakeextension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	"meganex/globals"
)

func CleanupSearchMatchmakeSession(matchmakeSession *matchmakingtypes.MatchmakeSession) {
	for _, index := range globals.NexConfig.MatchmakingZeroAttributes {
		matchmakeSession.Attributes[index] = 0
	}
}
func CleanupMatchmakeSessionSearchCriterias(searchCriterias types.List[matchmakingtypes.MatchmakeSessionSearchCriteria]) {
	for _, criteria := range searchCriterias {
		for _, index := range globals.NexConfig.MatchmakingZeroAttributes {
			criteria.Attribs[index] = ""
		}
	}
}

// GetPlayingSession is a stub impl until this can be added into nex-protocols-common-go for reals.
// https://github.com/PretendoNetwork/nex-protocols-common-go/issues/48
// Used by at least Splatoon, Yo-Kai Watch 2, and Tri Force Heroes
func GetPlayingSession(err error, packet nex.PacketInterface, callID uint32, _ types.List[types.PID]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	connection := packet.Sender().(*nex.PRUDPConnection)
	endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)

	// Empty list for stub
	lstPlayingSession := types.NewList[matchmakingtypes.PlayingSession]()

	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	lstPlayingSession.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = matchmakeextension.ProtocolID
	rmcResponse.MethodID = matchmakeextension.MethodGetPlayingSession
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
