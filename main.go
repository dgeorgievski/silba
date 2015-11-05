package main

import (
  "os"
  "silba/commands"
  "github.com/codegangsta/cli"
  "github.com/docker/machine/log"
)

func cmdNotFound(c *cli.Context, command string) {
	log.Fatalf(
		"%s: '%s' is not a %s command. See '%s --help'.\n",
		c.App.Name,
		command,
		c.App.Name,
		c.App.Name,
	)
}

var os_args []string

func main() {
  os_args = os.Args
  app := cli.NewApp()
  app.Name = "silba"
  app.Usage = "streaming metrics"
  app.Version = "1.0.0"

  app.Action = func (c *cli.Context) {
    log.Infof("No commands given. Using default.")
    commands.InitProbes(c)
  }

  app.Flags = []cli.Flag {
          cli.DurationFlag{
              EnvVar: "SILBA_INTERVAL",
              Name: "interval, i",
              Value: 30000000000, //3secs
              Usage: "Time interval between two consequtive samplings",
          },
          cli.IntFlag{
              EnvVar: "SILBA_COUNT",
              Name: "count, c",
              Value: 0,
              Usage: "Number of samplings",
          },
          cli.IntFlag{
              EnvVar: "SILBA_ZMQ_PORT",
              Name: "port, p",
              Value: 5555,
              Usage: "ZMQ pub port",
          },
          cli.StringFlag{
              EnvVar: "SILBA_PROBES",
              Name: "probes, P",
              Value: "cpu,mem,uptime",
              Usage: "List of comma delimited probes",
          },
        }

  app.Run(os_args)

}
