package commands

import (
	"cli/internal/parameters"
	"cli/internal/render"
	"context"

	"github.com/urfave/cli/v3"
)

func InitializeCommands() *cli.Command {
	return &cli.Command{
		Name:  "radio",
		Usage: "Radia Renderer",
		Commands: []*cli.Command{
			&cli.Command{
				Name:    "render",
				Usage:   "Render a Radia Scene",
				Aliases: []string{"r"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "scene",
						Usage:    "Scene Description JSON",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "width",
						Usage:    "Output width",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "height",
						Usage:    "Output height",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "samples",
						Usage:    "Number of samples",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "bounces",
						Usage:    "Bounce Limit",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "threads",
						Usage:    "Thread limit (Set 0 for unlimited)",
						Required: false,
						Value:    0,
					},
					&cli.StringFlag{
						Name:     "output",
						Usage:    "Output filename",
						Required: true,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					scene := cmd.String("scene")
					params := &parameters.RenderParameters{
						Width:   cmd.Int("width"),
						Height:  cmd.Int("height"),
						Samples: cmd.Int("samples"),
						Bounces: cmd.Int("bounces"),
						Threads: cmd.Int("threads"),
						Output:  cmd.String("output"),
					}
					return render.Scene(scene, params)
				},
			},
		},
	}
}
