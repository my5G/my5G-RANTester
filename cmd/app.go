package main

import (
	"my5G-RANTester/config"
	"my5G-RANTester/internal/templates"

	// "fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const version = "1.0.1"

func init() {

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	if cfg.Logs.Level == 0 {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.Level(cfg.Logs.Level))
	}

	spew.Config.Indent = "\t"

	log.Info("my5G-RANTester version " + version)
}

func main() {

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "ue",
				Aliases: []string{"ue"},
				Usage:   "Testing an ue attached with configuration",
				Action: func(c *cli.Context) error {
					name := "Testing an ue attached with configuration"
					cfg := config.Data

					log.Info("---------------------------------------")
					log.Info("[TESTER] Starting test function: ", name)
					log.Info("[TESTER][UE] Number of UEs: ", 1)
					log.Info("[TESTER][GNB] Control interface IP/Port: ", cfg.GNodeB.ControlIF.Ip, "/", cfg.GNodeB.ControlIF.Port)
					log.Info("[TESTER][GNB] Data interface IP/Port: ", cfg.GNodeB.DataIF.Ip, "/", cfg.GNodeB.DataIF.Port)
					log.Info("[TESTER][AMF] AMF IP/Port: ", cfg.AMF.Ip, "/", cfg.AMF.Port)
					log.Info("---------------------------------------")
					templates.TestAttachUeWithConfiguration()
					return nil
				},
			},
			{
				Name:    "non3gpp-ue",
				Aliases: []string{"ue"},
				Usage:   "Testing a non3gpp-ue attached with configuration",
				Action: func(c *cli.Context) error {
					name := "Testing an ue attached with configuration"
					cfg := config.Data

					log.Info("---------------------------------------")
					log.Info("[TESTER] Starting test function: ", name)
					log.Info("[TESTER][UE] Number of Non3GPP UEs: ", 1)
					log.Info("[TESTER][AMF] AMF IP/Port: ", cfg.AMF.Ip, "/", cfg.AMF.Port)
					log.Info("---------------------------------------")
					templates.TestAttachNon3gppUeWithConfiguration()
					return nil
				},
			},
			{
				Name:    "gnb",
				Aliases: []string{"gnb"},
				Usage:   "Testing an gnb attached with configuration",
				Action: func(c *cli.Context) error {
					name := "Testing an gnb attached with configuration"
					cfg := config.Data

					log.Info("---------------------------------------")
					log.Info("[TESTER] Starting test function: ", name)
					log.Info("[TESTER][GNB] Number of GNBs: ", 1)
					log.Info("[TESTER][GNB] Control interface IP/Port: ", cfg.GNodeB.ControlIF.Ip, "/", cfg.GNodeB.ControlIF.Port)
					log.Info("[TESTER][GNB] Data interface IP/Port: ", cfg.GNodeB.DataIF.Ip, "/", cfg.GNodeB.DataIF.Port)
					log.Info("[TESTER][AMF] AMF IP/Port: ", cfg.AMF.Ip, "/", cfg.AMF.Port)
					log.Info("---------------------------------------")
					templates.TestAttachGnbWithConfiguration()
					return nil
				},
			},
			{
				Name:    "load-test",
				Aliases: []string{"load-test"},
				Usage: "\nLoad endurance stress tests.\n" +
					"Example for testing multiple UEs: load-test -n 5 \n",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "number-of-ues", Value: 1, Aliases: []string{"n"}},
				},
				Action: func(c *cli.Context) error {
					var numUes int
					name := "Testing registration of multiple UEs"
					cfg := config.Data

					if c.IsSet("number-of-ues") {
						numUes = c.Int("number-of-ues")
					} else {
						log.Info(c.Command.Usage)
						return nil
					}

					log.Info("---------------------------------------")
					log.Info("[TESTER] Starting test function: ", name)
					log.Info("[TESTER][UE] Number of UEs: ", numUes)
					log.Info("[TESTER][GNB] gNodeB control interface IP/Port: ", cfg.GNodeB.ControlIF.Ip, "/", cfg.GNodeB.ControlIF.Port)
					log.Info("[TESTER][GNB] gNodeB data interface IP/Port: ", cfg.GNodeB.DataIF.Ip, "/", cfg.GNodeB.DataIF.Port)
					log.Info("[TESTER][AMF] AMF IP/Port: ", cfg.AMF.Ip, "/", cfg.AMF.Port)
					log.Info("---------------------------------------")
					templates.TestMultiUesInQueue(numUes)

					return nil
				},
			},
			{
				Name:    "amf-load-loop",
				Aliases: []string{"amf-load-loop"},
				Usage: "\nTest AMF responses in interval\n" +
					"Example for generating 20 requests to AMF per second in interval of 20 seconds: amf-load-loop -n 20 -t 20\n",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "number-of-requests", Value: 1, Aliases: []string{"n"}},
					&cli.IntFlag{Name: "time", Value: 1, Aliases: []string{"t"}},
				},
				Action: func(c *cli.Context) error {
					var time int
					var numRqs int

					name := "Test AMF responses in interval"
					cfg := config.Data

					numRqs = c.Int("number-of-requests")
					time = c.Int("time")

					log.Info("---------------------------------------")
					log.Warn("[TESTER] Starting test function: ", name)
					log.Warn("[TESTER][UE] Number of Requests per second: ", numRqs)
					log.Info("[TESTER][GNB] gNodeB control interface IP/Port: ", cfg.GNodeB.ControlIF.Ip, "/", cfg.GNodeB.ControlIF.Port)
					log.Info("[TESTER][GNB] gNodeB data interface IP/Port: ", cfg.GNodeB.DataIF.Ip, "/", cfg.GNodeB.DataIF.Port)
					log.Info("[TESTER][AMF] AMF IP/Port: ", cfg.AMF.Ip, "/", cfg.AMF.Port)
					log.Info("---------------------------------------")
					log.Warn("[TESTER][GNB] Total of AMF Responses in the interval:", templates.TestRqsLoop(numRqs, time))
					return nil
				},
			},
			{
				Name:    "ue-latency-interval",
				Aliases: []string{"ue-latency-interval"},
				Usage: "\nTesting UE latency in registration\n" +
					"Testing UE latency for 20 requests: ue-latency-interval -n 20\n",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "number-of-requests", Value: 1, Aliases: []string{"n"}},
				},
				Action: func(c *cli.Context) error {
					var requests int

					name := "Testing UE latency in registration"
					cfg := config.Data

					requests = c.Int("number-of-requests")

					log.Info("---------------------------------------")
					log.Warn("[TESTER] Starting test function: ", name)
					log.Warn("[TESTER][UE] Number of requests: ", requests)
					log.Info("[TESTER][GNB] Control interface IP/Port: ", cfg.GNodeB.ControlIF.Ip, "/", cfg.GNodeB.ControlIF.Port)
					log.Info("[TESTER][GNB] Data interface IP/Port: ", cfg.GNodeB.DataIF.Ip, "/", cfg.GNodeB.DataIF.Port)
					log.Info("[TESTER][AMF] AMF IP/Port: ", cfg.AMF.Ip, "/", cfg.AMF.Port)
					log.Info("---------------------------------------")
					log.Warn("[TESTER][UE] Average of the latency for a queue of requests: ", templates.TestUesLatencyInInterval(requests)/int64(requests), "ms")
					return nil
				},
			},
			{
				Name:    "Test availability of AMF",
				Aliases: []string{"amf-availability"},
				Usage: "\nTest availability of AMF in interval\n" +
					"Test availability of AMF in 20 seconds: amf-availability -t 20\n",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "time", Value: 1, Aliases: []string{"t"}},
				},
				Action: func(c *cli.Context) error {
					var time int

					name := "Test availability of AMF"
					cfg := config.Data

					time = c.Int("time")

					log.Info("---------------------------------------")
					log.Warn("[TESTER] Starting test function: ", name)
					log.Warn("[TESTER][UE] Interval of test: ", time, " seconds")
					log.Info("[TESTER][GNB] Control interface IP/Port: ", cfg.GNodeB.ControlIF.Ip, "/", cfg.GNodeB.ControlIF.Port)
					log.Info("[TESTER][GNB] Data interface IP/Port: ", cfg.GNodeB.DataIF.Ip, "/", cfg.GNodeB.DataIF.Port)
					log.Info("[TESTER][AMF] AMF IP/Port: ", cfg.AMF.Ip, "/", cfg.AMF.Port)
					log.Info("---------------------------------------")
					templates.TestAvailability(time)
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
