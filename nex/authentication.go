package nex

import (
	"fmt"
	"meganex/globals"
	"strings"

	"github.com/PretendoNetwork/nex-go/v2"
)

func StartAuthenticationServer() {
	globals.AuthenticationServer = nex.NewPRUDPServer()
	globals.AuthenticationServer.ByteStreamSettings.UseStructureHeader = globals.NexConfig.UseStructureHeader
	globals.AuthenticationServer.LibraryVersions.SetDefault(&globals.NexConfig.LibraryVersion)
	globals.AuthenticationServer.AccessKey = globals.NexConfig.AccessKey

	globals.AuthenticationEndpoint = nex.NewPRUDPEndPoint(1)
	globals.AuthenticationEndpoint.ServerAccount = globals.AuthenticationServerAccount
	globals.AuthenticationEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.AuthenticationEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.AuthenticationServer.BindPRUDPEndPoint(globals.AuthenticationEndpoint)

	globals.AuthenticationEndpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		header := fmt.Sprintf("== %v - Auth ==", globals.NexConfig.GameName)
		globals.Logger.Info(header)
		globals.Logger.Infof("Protocol ID: %d", request.ProtocolID)
		globals.Logger.Infof("Method ID: %d", request.MethodID)
		globals.Logger.Info(strings.Repeat("=", len(header)))
	})

	globals.AuthenticationEndpoint.OnError(func(err *nex.Error) {
		globals.Logger.Errorf("Auth: %v", err)
	})

	registerCommonAuthenticationServerProtocols()

	globals.AuthenticationServer.Listen(int(globals.NexConfig.AuthPort))
}
