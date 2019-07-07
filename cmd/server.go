package cmd

import (
	"os"
	"time"

	"github.com/gomodule/redigo/redis"

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
			rediscl, err := redis.Dial("tcp", Conf.CacheBind)
			if err != nil {
				log.Fatal().Err(err).Msg("")
			}
			defer rediscl.Close()

			go func() {
				for {
					if _, err := rediscl.Do("PING"); err != nil {
						log.Error().Err(err).Msg("redis ping")
					}
					time.Sleep(1 * time.Second)
				}
			}()

			cachesvc := pac.NewCache(rediscl)
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
