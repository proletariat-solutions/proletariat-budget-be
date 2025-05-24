package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"ghorkov32/proletariat-budget-be/config"
	"ghorkov32/proletariat-budget-be/internal/adapter/mongodb"
	"ghorkov32/proletariat-budget-be/internal/adapter/resthttp"
	"ghorkov32/proletariat-budget-be/internal/common"
	"ghorkov32/proletariat-budget-be/internal/core/usecase"
	"ghorkov32/proletariat-budget-be/openapi"
	"github.com/inv-cloud-platform/hub-com-auth-go/hubauth"
	"github.com/inv-cloud-platform/hub-com-tools-go/hubmiddlewares"
	"github.com/inv-cloud-platform/hub-com-tools-go/hubmongo"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("application starting ...")
	configs := config.Load()
	common.SetupLogger(configs.App.LogLevel)

	appCtx := context.Background()
	token := hubauth.NewTokenV2(hubauth.TokenConfig{
		Host:     configs.Auth.KeycloakHost,
		Realm:    configs.Auth.KeycloakRealm,
		Client:   configs.Auth.ClientID,
		Secret:   configs.Auth.ClientSecret,
		Username: configs.Auth.ServiceUsername,
		Password: configs.Auth.ServicePassword,
	})

	token.StartAutoRefreshV2(appCtx, func(err error) {
		log.Fatal().Err(err).Msg("unable to generate token")
	})

	// https://github.com/inv-cloud-platform/hub-com-tools-go/blob/main/hubmongo/README.md#connect
	mongoClient, errMongo := hubmongo.ConnectV2(hubmongo.DefaultConnection())
	if errMongo != nil {
		log.Fatal().Err(errMongo).Msg("unable to connect to mongodb")
	}

	scopeLookup, errLookup := hubauth.NewScopeLookup(hubauth.DefaultScopeOptions().WithAddress(configs.Auth.LookupApiHost))
	if errLookup != nil {
		log.Fatal().Err(errLookup).Msg("unable to create scope lookup")
	}

	accessChecker := usecase.NewAccessChecker(scopeLookup)
	userRepo := mongodb.NewUserRepo(mongoClient)
	userSvc := usecase.NewCreateUser(userRepo)
	userController := resthttp.NewUserController(userSvc, accessChecker)
	userHandler := openapi.NewStrictHandler(userController, nil)

	auth := hubauth.NewAuthorization(
		configs.Auth.KeycloakHost,
		configs.Auth.KeycloakRealm,
	)

	oapiSpec, errSw := openapi.GetSwagger()
	if errSw != nil {
		log.Warn().Err(errSw).Msg("unable to fetch openapi spec")
	}

	httpServer := resthttp.NewHTTPServer(
		configs.App,
		openapi.Handler(userHandler),
		resthttp.MetricsCollector,
		// TODO:  Remove if you don't need. Generally used by BFFs.
		hubmiddlewares.Proxy("/proxy/myApi", "http://localhost:8080", &hubmiddlewares.ProxyOptions{
			ProxyLogic: hubmiddlewares.PROXY_SELECTED,
			Endpoints: map[string][]string{
				"/something": {http.MethodGet},
			},
		}),
		hubmiddlewares.RequestTrack("/health"),
		hubmiddlewares.Health("/health",
			hubmiddlewares.HealthCheckMongoDB("mongodb", mongoClient),
		),
		hubmiddlewares.Spec(
			hubmiddlewares.SpecOptionsTitle(oapiSpec.Info.Title),
			hubmiddlewares.SpecOptionsObjectSpec(oapiSpec),
		),
		auth.Middleware(&hubauth.MiddlewareConfig{
			Strategy: hubauth.AllOf,
			Domains:  []string{common.Domain},
		}),
	)
	go httpServer.Start()
	defer func(ctx context.Context) {
		errS := httpServer.Shutdown(ctx)
		if errS != nil {
			log.Err(errS).Msg("http server shutdown error")
		}
	}(appCtx)

	ctx, stop := signal.NotifyContext(appCtx, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	log.Info().Msg("waiting for app exiting conditions")

	<-ctx.Done()
	log.Info().Msg("application stopping ...")
}
