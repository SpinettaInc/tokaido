package cmd

import (
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/tok"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/spf13/cobra"
)

// PsCmd - `tok ps`
var PsCmd = &cobra.Command{
	Use:   "ps",
	Short: "Show Tokaido containers and status",
	Long:  "Alias for running `docker-compose ps`",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.TokConfig("ps")
		utils.CheckCmdHard("docker-compose")

		docker.HardCheckTokCompose()

		tok.Ps()
	},
}
