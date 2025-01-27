package main

import (
	"database/sql"
	"fmt"
	pb "github.com/PretendoNetwork/grpc/go/account"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"meganex/globals"
	"os"
	"strconv"
	"strings"
)

func initConfig() error {
	return envconfig.Process(globals.EnvPrefix, &globals.NexConfig)
}

func init() {
	globals.Logger = plogger.NewLogger()

	var err error

	err = godotenv.Load()
	if err != nil {
		globals.Logger.Warning("Error loading .env file")
	}

	err = initConfig()
	if err != nil {
		globals.Logger.Errorf("NEX config wasn't valid! %s", err.Error())
		// saveDefaultConfig TODO
		os.Exit(1)
	}
	globals.Logger.Infof(globals.NexConfig.FormatToString(0))

	kerberosPassword := os.Getenv(globals.EnvPrefix + "_KERBEROS_PASSWORD")
	authenticationServerPort := os.Getenv(globals.EnvPrefix + "_AUTHENTICATION_SERVER_PORT")
	secureServerHost := os.Getenv(globals.EnvPrefix + "_SECURE_SERVER_HOST")
	secureServerPort := os.Getenv(globals.EnvPrefix + "_SECURE_SERVER_PORT")
	accountGRPCHost := os.Getenv(globals.EnvPrefix + "_ACCOUNT_GRPC_HOST")
	accountGRPCPort := os.Getenv(globals.EnvPrefix + "_ACCOUNT_GRPC_PORT")
	accountGRPCAPIKey := os.Getenv(globals.EnvPrefix + "_ACCOUNT_GRPC_API_KEY")

	if strings.TrimSpace(kerberosPassword) == "" {
		globals.Logger.Warningf(globals.EnvPrefix+"_KERBEROS_PASSWORD environment variable not set. Using default password: %q", globals.KerberosPassword)
	} else {
		globals.KerberosPassword = kerberosPassword
	}

	globals.InitAccounts()

	if strings.TrimSpace(authenticationServerPort) == "" {
		globals.Logger.Error(globals.EnvPrefix + "_AUTHENTICATION_SERVER_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(authenticationServerPort); err != nil {
		globals.Logger.Errorf(globals.EnvPrefix+"_AUTHENTICATION_SERVER_PORT is not a valid port. Expected 0-65535, got %s", authenticationServerPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf(globals.EnvPrefix+"_AUTHENTICATION_SERVER_PORT is not a valid port. Expected 0-65535, got %s", authenticationServerPort)
		os.Exit(0)
	}

	if strings.TrimSpace(secureServerHost) == "" {
		globals.Logger.Error(globals.EnvPrefix + "_SECURE_SERVER_HOST environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(secureServerPort) == "" {
		globals.Logger.Error(globals.EnvPrefix + "_SECURE_SERVER_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(secureServerPort); err != nil {
		globals.Logger.Errorf(globals.EnvPrefix+"_SECURE_SERVER_PORT is not a valid port. Expected 0-65535, got %s", secureServerPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf(globals.EnvPrefix+"_SECURE_SERVER_PORT is not a valid port. Expected 0-65535, got %s", secureServerPort)
		os.Exit(0)
	}

	if strings.TrimSpace(accountGRPCHost) == "" {
		globals.Logger.Error(globals.EnvPrefix + "_ACCOUNT_GRPC_HOST environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(accountGRPCPort) == "" {
		globals.Logger.Error(globals.EnvPrefix + "_ACCOUNT_GRPC_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(accountGRPCPort); err != nil {
		globals.Logger.Errorf(globals.EnvPrefix+"_ACCOUNT_GRPC_PORT is not a valid port. Expected 0-65535, got %s", accountGRPCPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf(globals.EnvPrefix+"_ACCOUNT_GRPC_PORT is not a valid port. Expected 0-65535, got %s", accountGRPCPort)
		os.Exit(0)
	}

	if strings.TrimSpace(accountGRPCAPIKey) == "" {
		globals.Logger.Warning("Insecure gRPC server detected. " + globals.EnvPrefix + "_ACCOUNT_GRPC_API_KEY environment variable not set")
	}

	globals.GRPCAccountClientConnection, err = grpc.NewClient(fmt.Sprintf("dns:%s:%s", accountGRPCHost, accountGRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		globals.Logger.Criticalf("Failed to connect to account gRPC server: %v", err)
		os.Exit(0)
	}

	globals.GRPCAccountClient = pb.NewAccountClient(globals.GRPCAccountClientConnection)
	globals.GRPCAccountCommonMetadata = metadata.Pairs(
		"X-API-Key", accountGRPCAPIKey,
	)
	globals.Postgres, err = sql.Open("postgres", os.Getenv(globals.EnvPrefix+"_POSTGRES_URI"))
	if err != nil {
		globals.Logger.Critical(err.Error())
	}
	globals.Logger.Success("Connected to Postgres!")
}
