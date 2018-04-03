package main

import (
	"github.com/mitchellh/go-homedir"
	cli "gopkg.in/urfave/cli.v1"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Code Screen Shooter"
	app.Usage = "[filepath]"

	var language string
	var out string
	var style string

	home, err := homedir.Dir()
	onError(err)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "lang, l",
			Value:       "auto",
			Usage:       "Language",
			Destination: &language,
		},
		cli.StringFlag{
			Name:        "output, o",
			Value:       home + "/code_screenshot.png",
			Usage:       "Output File",
			Destination: &out,
		},
		cli.StringFlag{
			Name:        "style, s",
			Value:       "monokai",
			Usage:       "Theme Style",
			Destination: &style,
		},
	}

	renderer := NewRenderer()

	app.Action = func(c *cli.Context) error {

		contents, err := ioutil.ReadFile(os.Args[1])
		onError(err)

		if language == "auto" {
			language = "match:" + os.Args[1]
		}

		renderer.ChangeStyle(style)
		rgba := renderer.Render(string(contents), language)
		f, err := os.Create(out)
		onError(err)

		png.Encode(f, rgba)

		return nil
	}

	err = app.Run(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

func onError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
