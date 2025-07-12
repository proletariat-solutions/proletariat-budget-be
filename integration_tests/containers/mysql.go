package containers

import (
	"context"
	"ghorkov32/proletariat-budget-be/config"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// define constants that will be used for the creation of the mysql container
const (
	mysqlImage = "mysql:8.0"
	mysqlPort  = "3306"
)

type MysqlContainer struct {
	container testcontainers.Container
}

func NewMysqlContainer() *MysqlContainer {
	return &MysqlContainer{}
}

// InitContainer initializes the container and return the config needed to connect with it
func (m *MysqlContainer) InitContainer(config *config.MySQL) (*config.MySQL, error) {
	log.Info().Msg("initializing mysql container...")
	// create the container request
	containerReq := testcontainers.ContainerRequest{
		Image:        mysqlImage,
		ExposedPorts: []string{mysqlPort + "/tcp"},
		Env: map[string]string{
			"MYSQL_USER":          config.User,
			"MYSQL_PASSWORD":      config.Password,
			"MYSQL_ROOT_PASSWORD": config.Password,
			"MYSQL_DATABASE":      config.Database,
			// you can add more environment variables here
		},
		SkipReaper: true,
		WaitingFor: wait.ForListeningPort(mysqlPort + "/tcp"), // will wait for the container to be up and running
	}

	ctx := context.Background()

	// creates a new container
	mySQLContainer, err := testcontainers.GenericContainer(
		ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		},
	)

	if err != nil {
		return nil, err
	}

	host, _ := mySQLContainer.Host(ctx)
	config.Host = host

	//  testcontainers map his port to an random external port, so we grab the port with this method
	port, _ := mySQLContainer.MappedPort(ctx, mysqlPort+"/tcp")
	portNumber := port.Int()

	m.container = mySQLContainer

	// create a new mongoConfig with testcontainer information
	config.Port = strconv.Itoa(portNumber)

	log.Info().Interface("mysql container initialized successfully", *config)

	return config, nil
}

// DestroyContainer destroys the container.
func (m *MysqlContainer) DestroyContainer() error {
	return m.container.Terminate(context.Background())
}
