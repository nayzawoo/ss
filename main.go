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
	app.Name = "Code Screenshot"
	app.Usage = "ss <options> [sourcecode]"

	home, err := homedir.Dir()
	checkError(err)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lang, l",
			Value: "auto",
			Usage: "`LANGUAGE`",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "home",
			Usage: "Output `FILE`",
		},
		cli.StringFlag{
			Name:  "style, s",
			Value: "monokai",
			Usage: "Theme `STYLE`",
		},
	}

	renderer := NewRenderer()

	app.Action = func(c *cli.Context) error {
		if len(os.Args) == 1 {
			cli.ShowAppHelp(c)
			os.Exit(0)
		}

		lang := c.String("lang")
		style := c.String("style")
		output := c.String("output")

		filePath := c.Args().Get(0)
		contents, err := ioutil.ReadFile(filePath)
		checkError(err)

		if lang == "auto" {
			lang = "match:" + filePath
		}

		renderer.ChangeStyle(style)
		rgba := renderer.Render(string(contents), lang)

		if output == "home" {
			output = home + "/code_screenshot.png"
		}

		f, err := os.Create(output)
		checkError(err)

		png.Encode(f, rgba)

		return nil
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
