package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"otsukai"
	"otsukai/parser"
	"otsukai/runtime"
	"otsukai/runtime/context"
)

func run(c *cli.Context) error {
	recipe := c.String("recipe")

	content, err := os.ReadFile(recipe)
	if err != nil {
		log.Fatal(err)
		return err
	}

	ruby, err := parser.Parser.ParseString("", string(content)+"\n")
	if err != nil {
		fmt.Println(err)
		return err
	}

	ctx := context.NewContext(ruby)
	if err := runtime.Run(ctx); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func test(c *cli.Context) error {
	recipe := c.String("recipe")
	otsukai.Infof("check syntax for %s", recipe)

	content, err := os.ReadFile(recipe)
	if err != nil {
		otsukai.Fatalf("failed to read recipe: %s", err.Error())
		return err
	}

	_, err = parser.Parser.ParseString("", string(content)+"\n")
	if err != nil {
		otsukai.Errf("check syntas for %s is failed: %s", recipe, err.Error())
		return err
	}

	otsukai.Successf("check syntax for %s is success", recipe)
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
