package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"

	"github.com/vbatoufflet/goverview/pkg/logger"
	"github.com/vbatoufflet/goverview/pkg/yaml"
)

const (
	defaultConfigFile string = "/etc/goverview/goverview.yml"
)

var (
	flagConfigFile string
	flagHelp       bool
	flagVersion    bool
	version        string
)

func main() {
	var (
		c   Config
		s   *Server
		sc  chan os.Signal
		err error
	)

	flag.StringVar(&flagConfigFile, "c", defaultConfigFile, "configuration file path")
	flag.BoolVar(&flagHelp, "h", false, "display this help and exit")
	flag.BoolVar(&flagVersion, "V", false, "display software version and exit")
	flag.Usage = func() { printUsage(os.Stderr); os.Exit(1) }
	flag.Parse()

	if flagHelp {
		printUsage(os.Stdout)
		os.Exit(0)
	} else if flagVersion {
		fmt.Printf("%s version %s\nGo version: %s (%s)\n", path.Base(os.Args[0]), version, runtime.Version(),
			runtime.Compiler)
		os.Exit(0)
	} else if flagConfigFile == "" {
		fmt.Fprintf(os.Stderr, "Error: configuration file path is mandatory\n")
		printUsage(os.Stderr)
		os.Exit(1)
	}

	// Load server configuration
	c = Config{}
	if err = yaml.Load(flagConfigFile, &c); err != nil {
		goto end
	}

	// Set log output
	if err = logger.Init(c.LogPath, c.LogLevel); err != nil {
		goto end
	}

	// Start server instance
	s = NewServer(c)

	// Handle server signals
	sc = make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	go func() {
		for sig := range sc {
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				s.Stop()

			case syscall.SIGUSR1:
				s.Refresh()
			}
		}
	}()

	err = s.Run()

end:
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func printUsage(fd io.Writer) {
	fmt.Fprintf(fd, "%s [OPTIONS] -c <path>", path.Base(os.Args[0]))
	fmt.Fprint(fd, "\n\nOptions:\n")

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(fd, "   -%s  %s\n", f.Name, f.Usage)
	})
}
