package cmd

import (
	"os"

	"github.com/mediocregopher/radix/v3"

	"github.com/rustwizard/tpa/internal/cache"

	"github.com/rustwizard/tpa/internal/pac"

	"github.com/rustwizard/tpa/internal/server/transport/http"

	"github.com/rustwizard/tpa/internal/server"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "the command runs autocomplete server",

	Run: func(cmd *cobra.Command, args []string) {
		log := zerolog.New(os.Stdout).With().Caller().Logger().With().Str("pkg", "server").Logger()

		log.Debug().Msgf("config: %v", Conf)

		switch Conf.Transport {
		case "http":
			rad, err := radix.Dial("tcp", Conf.CacheBind)
			if err != nil {
				log.Fatal().Err(err).Msg("")
			}
			defer rad.Close()

			radcl := radix.Client(rad)
			cachesvc := cache.NewService(radcl)
			pacsvc := pac.NewService(log, Conf.RemoteAPIPath, cachesvc)
			handler := http.NewHandler(log, pacsvc)
			srv := http.NewServer(log, &Conf, handler)
			if err := srv.Run(); err != nil {
				log.Fatal().Err(err).Msg("")
			}
		default:
			log.Fatal().Err(server.ErrTransport).Msg("")
		}

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVar(&Conf.Bind, "bind", "", "addr:port")
	serverCmd.PersistentFlags().DurationVar(&Conf.RequestTTL, "requesttl", 0, "provide request ttl")
	serverCmd.PersistentFlags().StringVar(&Conf.Transport, "transport", "", "http")
}
