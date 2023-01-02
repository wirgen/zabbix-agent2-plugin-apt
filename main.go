package main

import (
	"flag"
	"fmt"
	"git.zabbix.com/ap/apt/plugin"
	"git.zabbix.com/ap/plugin-support/plugin/comms"
	"git.zabbix.com/ap/plugin-support/plugin/container"
	"os"
)

const pluginVersion = 0

func main() {
	handleFlags()

	h, err := container.NewHandler(plugin.Impl.Name())
	if err != nil {
		panic(fmt.Sprintf("failed to create plugin handler %s", err.Error()))
	}
	plugin.Impl.Logger = &h

	err = h.Execute()
	if err != nil {
		panic(fmt.Sprintf("failed to execute plugin handler %s", err.Error()))
	}
}

func handleFlags() {
	var versionFlag bool
	const (
		versionDefault     = false
		versionDescription = "Print program version and exit"
	)
	flag.BoolVar(&versionFlag, "version", versionDefault, versionDescription)
	flag.BoolVar(&versionFlag, "V", versionDefault, versionDescription+" (shorthand)")

	var helpFlag bool
	const (
		helpDefault     = false
		helpDescription = "Display this help message"
	)
	flag.BoolVar(&helpFlag, "help", helpDefault, helpDescription)
	flag.BoolVar(&helpFlag, "h", helpDefault, helpDescription+" (shorthand)")

	flag.Parse()

	if helpFlag || len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "V", "version":
			comms.PrintVersion(plugin.Impl.Name(), copyrightMessage(), pluginVersion)
			os.Exit(0)
		}
	})
}

func copyrightMessage() string {
	return ""
}
