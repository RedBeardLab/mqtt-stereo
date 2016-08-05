package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/siscia/mqtt-player/backend"
)

func main() {
	app := cli.NewApp()
	app.Name = "mqtt-player"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "record",
			Usage: "What file use to record and playback",
			Value: "mqtt-record.txt",
		},
		cli.StringFlag{
			Name:  "topic",
			Usage: "What topic listen and what topic play back",
			Value: "/#",
		},
		cli.StringFlag{
			Name:  "url",
			Usage: "Where to listen to the MQTT broker",
			Value: "localhost",
		},
		cli.IntFlag{
			Name:  "port",
			Usage: "What port to use to connect to the broker",
			Value: 1883,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "record",
			Aliases: []string{"r", "rec"},
			Usage:   "Record traffic from a MQTT broker",
			Action: func(c *cli.Context) error {
				fmt.Println(c.GlobalString("url"))
				backend.StartRecording(c)
				return nil
			},
		},
		{
			Name:    "play",
			Aliases: []string{"p"},
			Usage:   "Play back previous registered traffic to a MQTT broker",
			Action: func(c *cli.Context) error {
				fmt.Println("Playing back the recorded traffic from: ", c.GlobalString("record"))
				backend.PlayBack(c)
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "ff",
					Usage: "Fast forward, if true plays all the messages in the file without respecting the times differences"},

				cli.BoolFlag{
					Name:  "loop",
					Usage: "Loop the player and keep playing all the messages, in order, indefinitely",
				},
			},
		},
	}

	app.Run(os.Args)
}
