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
	app.Commands = []cli.Command{
		cli.Command{
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
					Value: 10,
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
				return CreatePuzzle(input, output, c.Int("size"), time.Now().UTC().UnixNano())
			},
		},
	}

	app.Run(os.Args)
}
