package nex

import (
	"github.com/PretendoNetwork/nex-go/v2/types"
	matchmakingtypes "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
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
