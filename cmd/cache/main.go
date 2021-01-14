package main

import (
	"fmt"
	"github.com/KyberNetwork/cache/fetcher"
	persister "github.com/KyberNetwork/cache/persister"
	"log"
	"os"
	"runtime"
	"sort"

	cli "gopkg.in/urfave/cli.v1"
)

type fetcherFunc func(memPersister persister.MemoryPersister, fetcher *fetcher.Fetcher)

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	//set log for server
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app := cli.NewApp()
	app.Name = "Kyber Swap Cache"
	app.Usage = "Cache"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{}

	app.Commands = []cli.Command{}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Action = cmdMain

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func cmdMain(ctx *cli.Context) error {
	core, err := NewCacheCore()
	if err != nil {
		log.Println(err)
		return err
	}
	return core.Run()
}
