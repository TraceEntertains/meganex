package nex

import (
	"fmt"
	"github.com/PretendoNetwork/nex-go/v2"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"meganex/globals"
	"strings"
)

func StartSecureServer() {
	globals.SecureServer = nex.NewPRUDPServer()
	globals.SecureServer.ByteStreamSettings.UseStructureHeader = globals.NexConfig.UseStructureHeader
	globals.SecureServer.LibraryVersions.SetDefault(&globals.NexConfig.LibraryVersion)
	globals.SecureServer.AccessKey = globals.NexConfig.AccessKey

	globals.SecureEndpoint = nex.NewPRUDPEndPoint(1)
	globals.SecureEndpoint.IsSecureEndPoint = true
	globals.SecureEndpoint.ServerAccount = globals.SecureServerAccount
	globals.SecureEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.SecureEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.SecureServer.BindPRUDPEndPoint(globals.SecureEndpoint)

	globals.SecureEndpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		header := fmt.Sprintf("==%v - Secure==", globals.NexConfig.GameName)
		globals.Logger.Info(header)
		globals.Logger.Infof("Protocol ID: %d", request.ProtocolID)
		globals.Logger.Infof("Method ID: %d", request.MethodID)
		globals.Logger.Info(strings.Repeat("=", len(header)))
	})

	globals.SecureEndpoint.OnError(func(err *nex.Error) {
		globals.Logger.Errorf("Secure: %v", err)
	})

	globals.MatchmakingManager = common_globals.NewMatchmakingManager(globals.SecureEndpoint, globals.Postgres)

	registerCommonSecureServerProtocols()

	globals.SecureServer.Listen(int(globals.NexConfig.SecurePort))
}
