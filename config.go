package ocpi

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/mikhailbolshakov/kit/cluster"
	kitConfig "github.com/mikhailbolshakov/kit/config"
	"github.com/mikhailbolshakov/kit/grpc"
	kitHttp "github.com/mikhailbolshakov/kit/http"
	"github.com/mikhailbolshakov/kit/kafka"
	"github.com/mikhailbolshakov/kit/monitoring"
	"github.com/mikhailbolshakov/kit/profile"
	"github.com/mikhailbolshakov/kit/storages/pg"
	"os"
	"path/filepath"
)

var (
	// Meta service meta info
	Meta = cluster.NewMetaInfo("ocpi", kit.UUID(4))
	// Logger service logger
	Logger = kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel, Format: kit.FormatterJson})
)

func LF() kit.CLoggerFunc {
	return func() kit.CLogger {
		return kit.L(Logger).Srv(Meta.ServiceCode()).Node(Meta.InstanceId())
	}
}

func L() kit.CLogger {
	return LF()()
}

type CfgStorages struct {
	Database *pg.DbClusterConfig
}

type CfgAdapter struct {
	Grpc *grpc.ClientConfig
}

type CfgOcpiPlatform struct {
	Id       string
	Name     string
	TokenA   string `config:"token-a"`
	Role     string
	Versions []string
}

type CfgOcpiParty struct {
	PartyId     string `config:"party-id"`
	CountryCode string `config:"country-code"`
	Roles       string
}

type CfgOcpiEmulator struct {
	Id        string
	Name      string
	TokenA    string `config:"token-a"`
	Role      string
	VersionEp string `config:"version-ep"`
	Url       string
	ApiKey    string `config:"api-key"`
}

type CfgWebHook struct {
	Mock    bool
	Timeout *int
}

type CfgOcpiLocal struct {
	Url      string
	ApiKey   string `config:"api-key"`
	Platform *CfgOcpiPlatform
	Party    *CfgOcpiParty
	Webhook  *CfgWebHook
}

type CfgOcpiRemote struct {
	Mock    bool
	Timeout *int
}

type CfgOcpiConfig struct {
	Local    *CfgOcpiLocal
	Remote   *CfgOcpiRemote
	Emulator *CfgOcpiEmulator
}

type Tests struct {
	WebhookUrl string `config:"webhook-url"`
}

type Config struct {
	Grpc       *grpc.ServerConfig
	Storages   *CfgStorages
	Log        *kit.LogConfig
	Monitoring *monitoring.Config
	Adapters   map[string]*CfgAdapter
	Ws         *kitHttp.Config
	Profile    *profile.Config
	Kafka      *kafka.BrokerConfig
	Http       *kitHttp.Config
	Ocpi       *CfgOcpiConfig
	Tests      *Tests
}

func LoadConfig() (*Config, error) {

	// get root folder from env
	rootPath := os.Getenv("ROOT")
	if rootPath == "" {
		return nil, kitConfig.ErrEnvRootPathNotSet("ROOT")
	}

	// config path
	configPath := filepath.Join(rootPath, Meta.ServiceCode(), "config.yml")

	// .env path
	envPath := filepath.Join(rootPath, Meta.ServiceCode(), ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		envPath = ""
	}

	// load config
	config := &Config{}
	err := kitConfig.NewConfigLoader(LF()).
		WithConfigPath(configPath).
		WithEnvPath(envPath).
		Load(config)

	if err != nil {
		return nil, err
	}
	return config, nil
}
