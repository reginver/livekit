package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/livekit/livekit-server/pkg/config"
	"github.com/livekit/livekit-server/pkg/logger"
	"github.com/livekit/livekit-server/pkg/server"
	"github.com/livekit/protocol/livekit"
)

var (
	// Version is set at build time via ldflags
	Version = "dev"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app := &cli.App{
		Name:    "livekit-server",
		Usage:   "LiveKit SFU server",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "path to LiveKit config file",
				EnvVars: []string{"LIVEKIT_CONFIG_FILE"},
			},
			&cli.StringFlag{
				Name:    "config-body",
				Usage:   "LiveKit config in YAML, read from stdin or env",
				EnvVars: []string{"LIVEKIT_CONFIG"},
			},
			&cli.StringFlag{
				Name:    "key-file",
				Usage:   "path to file that contains API keys/secrets",
				EnvVars: []string{"LIVEKIT_KEYS_FILE"},
			},
			&cli.StringFlag{
				Name:    "keys",
				Usage:   "api keys (key: secret\nkey2: secret2)",
				EnvVars: []string{"LIVEKIT_KEYS"},
			},
			&cli.StringFlag{
				Name:    "node-ip",
				Usage:   "IP address of the current node, used to advertise to other nodes",
				EnvVars: []string{"NODE_IP"},
			},
			&cli.StringFlag{
				Name:    "redis",
				Usage:   "Redis URL, used for distributed deployments",
				EnvVars: []string{"REDIS_URL"},
			},
			&cli.BoolFlag{
				Name:  "dev",
				Usage: "enable development mode (insecure, no TLS required)",
			},
		},
		Action: startServer,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startServer(c *cli.Context) error {
	conf, err := config.NewConfig(c.String("config"), c.String("config-body"), c)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger.InitFromConfig(conf.Logging)

	if conf.Development {
		logger.GetLogger().Infow("starting in development mode")
	}

	_ = livekit.NodeID(fmt.Sprintf("nd_%s", newNodeID()))

	s, err := server.NewLiveKitServer(conf)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	if err := s.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	logger.GetLogger().Infow("server started, press Ctrl+C to stop")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.GetLogger().Infow("shutting down server")
	s.Stop(false)
	return nil
}

func newNodeID() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
