package main

import (
	"errors"
	"github.com/mika-sandbox/otsukai/logger"
	"github.com/mika-sandbox/otsukai/parser"
	"github.com/mika-sandbox/otsukai/runtime"
	"github.com/mika-sandbox/otsukai/runtime/context"
	re "github.com/mika-sandbox/otsukai/runtime/errors"
	"github.com/urfave/cli/v2"
	"os"
)

func run(c *cli.Context) error {
	recipe := c.String("recipe")

	content, err := os.ReadFile(recipe)
	if err != nil {
		logger.Fatalf("failed to read recipe: %s", err)
		return re.RUNTIME_ERROR
	}

	ruby, err := parser.Parser.ParseString("", string(content)+"\n")
	if err != nil {
		logger.Errf("invalid syntax: %s", err)
		return re.SYNTAX_ERROR
	}

	ctx := context.NewContext(ruby)
	if err := runtime.Run(&ctx); err != nil {
		if errors.Is(err, re.EXECUTION_ERROR) || errors.Is(err, re.SYNTAX_ERROR) || errors.Is(err, re.RUNTIME_ERROR) {
			return err
		}

		logger.Fatalf("%s", err)
		return err
	}

	return nil
}

func test(c *cli.Context) error {
	recipe := c.String("recipe")
	logger.Infof("check syntax for %s", recipe)

	content, err := os.ReadFile(recipe)
	if err != nil {
		logger.Fatalf("failed to read recipe: %s", err.Error())
		return err
	}

	_, err = parser.Parser.ParseString("", string(content)+"\n")
	if err != nil {
		logger.Errf("check syntas for %s is failed: %s", recipe, err.Error())
		return err
	}

	logger.Successf("check syntax for %s is success", recipe)
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
		os.Exit(1)
	}
}
