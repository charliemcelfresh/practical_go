package cmd

import (
	"charliemcelfresh/practical_go/internal/jwt_maker"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateJwtCmd)
}

var generateJwtCmd = &cobra.Command{
	Use: "generate_jwt",
	Run: func(cmd *cobra.Command, args []string) {
		// args[0] == duration
		// args[1] == issuing service
		// args[2] == audience, ie user, admin, or service-to-service
		// args[3] == user_id or admin_id
		// eg
		// generate_jwt 1h some-service service-to-service 1
		// generate_jwt 1h some-service user 1
		// generate_jwt 1h some-service admin 1
		adminOrUserId, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			logrus.Fatal(err)
		}
		jwt_maker.Run(args[0], args[1], args[2], adminOrUserId)
	},
}
