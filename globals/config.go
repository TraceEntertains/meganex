package globals

import (
	"fmt"
	"github.com/PretendoNetwork/nex-go/v2"
	"strings"
)

const EnvPrefix string = "PN_MEGANEX"

type NexConfigSpec struct {
	GameName           string `required:"true"`
	UseStructureHeader bool
	LibraryVersion     nex.LibraryVersion
	AccessKey          string `required:"true"`

	//AuthHost string `default:"127.0.0.1" envconfig:"AUTHENTICATION_SERVER_HOST"`
	AuthPort   uint16 `default:"60000" envconfig:"AUTHENTICATION_SERVER_PORT"`
	SecureHost string `default:"127.0.0.1" envconfig:"SECURE_SERVER_HOST"`
	SecurePort uint16 `default:"60001" envconfig:"SECURE_SERVER_PORT"`

	SecureProtocols []string

	MatchmakingZeroAttributes []int
}

// FormatToString pretty-prints the NexConfigSpec using the provided indentation level
func (nc *NexConfigSpec) FormatToString(indentationLevel int) string {
	indentationValues := strings.Repeat("\t", indentationLevel+1)
	indentationEnd := strings.Repeat("\t", indentationLevel)

	var b strings.Builder

	b.WriteString("NexConfigSpec{\n")
	b.WriteString(fmt.Sprintf("%sGameName: %v,\n", indentationValues, nc.GameName))
	b.WriteString(fmt.Sprintf("%sUseStructureHeader: %v,\n", indentationValues, nc.UseStructureHeader))
	b.WriteString(fmt.Sprintf("%sLibraryVersion: %v.%v.%v-%s,\n", indentationValues, nc.LibraryVersion.Major, nc.LibraryVersion.Minor, nc.LibraryVersion.Patch, nc.LibraryVersion.GameSpecificPatch))
	b.WriteString(fmt.Sprintf("%sAccessKey: %v,\n", indentationValues, nc.AccessKey))
	//b.WriteString(fmt.Sprintf("%sAuthHost: %v,\n", indentationValues, nc.AuthHost))
	b.WriteString(fmt.Sprintf("%sAuthPort: %v,\n", indentationValues, nc.AuthPort))
	b.WriteString(fmt.Sprintf("%sSecureHost: %v,\n", indentationValues, nc.SecureHost))
	b.WriteString(fmt.Sprintf("%sSecurePort: %v,\n", indentationValues, nc.SecurePort))
	b.WriteString(fmt.Sprintf("%sSecureProtocols: %v,\n", indentationValues, nc.SecureProtocols))
	b.WriteString(fmt.Sprintf("%sMatchmakingZeroAttributes: %v,\n", indentationValues, nc.MatchmakingZeroAttributes))
	b.WriteString(fmt.Sprintf("%s}", indentationEnd))

	return b.String()
}

var NexConfig NexConfigSpec
