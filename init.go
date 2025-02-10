package main

import (
	"database/sql"
	"fmt"
	"meganex/globals"
	"os"
	"strconv"
	"strings"

	pbaccount "github.com/PretendoNetwork/grpc/go/account"
	pbfriends "github.com/PretendoNetwork/grpc/go/friends"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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
	friendsGRPCHost := os.Getenv(globals.EnvPrefix + "_FRIENDS_GRPC_HOST")
	friendsGRPCPort := os.Getenv(globals.EnvPrefix + "_FRIENDS_GRPC_PORT")
	friendsGRPCAPIKey := os.Getenv(globals.EnvPrefix + "_FRIENDS_GRPC_API_KEY")
	s3Endpoint := os.Getenv(globals.EnvPrefix + "_CONFIG_S3_ENDPOINT")
	s3AccessKey := os.Getenv(globals.EnvPrefix + "_CONFIG_S3_ACCESS_KEY")
	s3AccessSecret := os.Getenv(globals.EnvPrefix + "_CONFIG_S3_ACCESS_SECRET")
	s3SecureEnv := os.Getenv(globals.EnvPrefix + "_CONFIG_S3_SECURE")

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

	globals.GRPCAccountClient = pbaccount.NewAccountClient(globals.GRPCAccountClientConnection)
	globals.GRPCAccountCommonMetadata = metadata.Pairs(
		"X-API-Key", accountGRPCAPIKey,
	)

	if strings.TrimSpace(friendsGRPCHost) == "" {
		globals.Logger.Error(globals.EnvPrefix + "_FRIENDS_GRPC_HOST environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(friendsGRPCPort) == "" {
		globals.Logger.Error(globals.EnvPrefix + "_FRIENDS_GRPC_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(friendsGRPCPort); err != nil {
		globals.Logger.Errorf(globals.EnvPrefix+"_FRIENDS_GRPC_PORT is not a valid port. Expected 0-65535, got %s", friendsGRPCPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf(globals.EnvPrefix+"_FRIENDS_GRPC_PORT is not a valid port. Expected 0-65535, got %s", friendsGRPCPort)
		os.Exit(0)
	}

	if strings.TrimSpace(friendsGRPCAPIKey) == "" {
		globals.Logger.Warning("Insecure gRPC server detected. " + globals.EnvPrefix + "_FRIENDS_GRPC_API_KEY environment variable not set")
	}

	globals.GRPCFriendsClientConnection, err = grpc.NewClient(fmt.Sprintf("dns:%s:%s", friendsGRPCHost, friendsGRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		globals.Logger.Criticalf("Failed to connect to friends gRPC server: %v", err)
		os.Exit(0)
	}

	globals.GRPCFriendsClient = pbfriends.NewFriendsClient(globals.GRPCFriendsClientConnection)
	globals.GRPCFriendsCommonMetadata = metadata.Pairs(
		"X-API-Key", friendsGRPCAPIKey,
	)

	globals.Postgres, err = sql.Open("postgres", os.Getenv(globals.EnvPrefix+"_POSTGRES_URI"))
	if err != nil {
		globals.Logger.Critical(err.Error())
	}
	globals.Logger.Success("Connected to Postgres!")

	staticCredentials := credentials.NewStaticV4(s3AccessKey, s3AccessSecret, "")
	s3Secure, err := strconv.ParseBool(s3SecureEnv)
	if err != nil {
		globals.Logger.Warningf(globals.EnvPrefix+"_CONFIG_S3_SECURE environment variable not set. Using default value: %t", true)
		s3Secure = true
	}

	minIOClient, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  staticCredentials,
		Secure: s3Secure,
	})
	if err != nil {
		globals.Logger.Warningf("Failed to connect to S3: %v", err)
		globals.Logger.Warning("Datastore uploads may not work.")
	} else {
		globals.MinIOClient = minIOClient
	}
}
