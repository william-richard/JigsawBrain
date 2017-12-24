package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "generate"
	app.Usage = "Make a jigsaw puzzle"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Show debug info",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "generate",
			Usage: "Generate puzzles",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Usage: "input image to generate a puzzle from",
				},
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Direction for output puzzle",
				},
				cli.IntFlag{
					Name:  "size, s",
					Usage: "Piece size",
					Value: 50,
				},
			},
			Action: func(c *cli.Context) error {
				input := c.String("input")
				output := c.String("output")
				if len(input) == 0 {
					log.Fatal("Input is required.")
				}
				if len(output) == 0 {
					log.Fatal("Output is required")
				}
				size := c.Int("size")
				llog := log.WithFields(log.Fields{"input": input, "output": output, "size": size})
				llog.Debug("Starting to create puzzle")
				puzzle, err := CreatePuzzleFromFile(input, size, time.Now().UTC().UnixNano())
				if err != nil {
					llog.WithError(err).Error("Error creating puzzle from files")
					return err
				}

				llog.Debug("Loaded puzzle")
				err = puzzle.WriteToDirectory(output)
				if err != nil {
					llog.WithError(err).Error("Error writing puzzle to directory")
					return err
				}
				llog.Debug("Wrote puzzle to dir")

				return nil

			},
		},
	}

	app.Run(os.Args)
}
