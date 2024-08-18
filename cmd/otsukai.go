package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"otsukai"
)

func run(c *cli.Context) error {
	recipe := c.String("recipe")
	content, err := os.ReadFile(recipe)
	if err != nil {
		log.Fatal(err)
		return err
	}

	ruby, err := otsukai.Parser.ParseString("", string(content)+"\n")
	if err != nil {
		fmt.Println(err)
		return err
	}

	ctx := otsukai.NewContext(ruby)
	if err := ctx.Run(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func test(c *cli.Context) error {
	recipe := c.String("recipe")
	content, err := os.ReadFile(recipe)
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = otsukai.Parser.ParseString("", string(content))
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("syntax ok")
	return nil
}

func main() {
	app := &cli.App{
		Name:  "otsukai",
		Usage: "otsukai",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "run otsukai",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "recipe",
						Value: "otsukai.rb",
						Usage: "path for recipe",
					},
				},
				Action: run,
			},
			{
				Name:  "test",
				Usage: "test recipe syntax",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "recipe",
						Value: "otsukai.rb",
						Usage: "path for recipe",
					},
				},
				Action: test,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
