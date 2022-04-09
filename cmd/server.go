package cmd

import (
	"charliemcelfresh/practical_go/internal/admin_hooks"
	"charliemcelfresh/practical_go/internal/middlewares"
	practical_go2 "charliemcelfresh/practical_go/internal/practical_go"
	"charliemcelfresh/practical_go/internal/service_to_service_hooks"
	"charliemcelfresh/practical_go/internal/user_hooks"
	"charliemcelfresh/practical_go/rpc/practical_go"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/sirupsen/logrus"

	"github.com/twitchtv/twirp"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	rootCmd.AddCommand(severCmd)

}

var severCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		Run()
	},
}

func Run() {
	provider := practical_go2.NewProvider()

	mux := http.NewServeMux()

	userChainHooks := twirp.ChainHooks(
		user_hooks.Auth(),
		user_hooks.Logging(),
	)

	serviceToServiceChainHooks := twirp.ChainHooks(
		service_to_service_hooks.Auth(),
		service_to_service_hooks.Logging(),
	)

	adminChainHooks := twirp.ChainHooks(
		admin_hooks.Auth(),
		admin_hooks.Audit(),
		admin_hooks.Logging(),
	)

	// http(s)://<host>:/v1/user/practical_go.PracticalGo/CreateItem
	// http(s)://<host>:/v1/user/practical_go.PracticalGo/GetItem
	userHandler := practical_go.NewPracticalGoServer(provider, twirp.WithServerPathPrefix("/v1/user"), userChainHooks)
	mux.Handle(
		userHandler.PathPrefix(), middlewares.AddRequestBodyToContext(
			middlewares.AddJwtTokenToContext(
				userHandler,
			),
		),
	)

	// http(s)://<host>:/v1/admin/practical_go.PracticalGo/CreateItem
	// http(s)://<host>:/v1/admin/practical_go.PracticalGo/GetItem
	adminHandler := practical_go.NewPracticalGoServer(
		provider, twirp.WithServerPathPrefix("/v1/admin"), adminChainHooks,
	)
	mux.Handle(
		adminHandler.PathPrefix(), middlewares.AddRequestBodyToContext(
			middlewares.AddJwtTokenToContext(
				adminHandler,
			),
		),
	)

	// http(s)://<host>:/v1/internal/practical_go.PracticalGo/CreateItem
	// http(s)://<host>:/v1/internal/practical_go.PracticalGo/GetItem
	serviceToServiceHandler := practical_go.NewPracticalGoServer(
		provider, twirp.WithServerPathPrefix("/v1/internal"),
		serviceToServiceChainHooks,
	)
	mux.Handle(
		serviceToServiceHandler.PathPrefix(), middlewares.AddRequestBodyToContext(
			middlewares.AddJwtTokenToContext(
				serviceToServiceHandler,
			),
		),
	)

	http.ListenAndServe(":8080", mux)
}
