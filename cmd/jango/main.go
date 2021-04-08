package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/dinalt/jango"
	"github.com/dinalt/jango/transport"
	"github.com/dinalt/jango/videoroom"
)

const (
	defaultBaseURL     = "http://localhost"
	defaultAdminSecret = "janusoverlord"
	defaultAdminPort   = 7088
)

var (
	adminEndpoint string
	adminSecret   string
)

func init() {
	base := os.Getenv("JANGO_JANUS_URL_BASE")
	if base == "" {
		base = defaultBaseURL
	}

	pref := os.Getenv("JANGO_JANUS_PATH_PREFIX")

	adminPort, _ := strconv.Atoi(os.Getenv("JANGO_JANUS_ADMIN_PORT"))
	if adminPort <= 0 {
		adminPort = defaultAdminPort
	}

	var ok bool
	adminSecret, ok = os.LookupEnv("JANGO_JANUS_ADMIN_SECRET")
	if !ok {
		adminSecret = defaultAdminSecret
	}

	flag.StringVar(&base, "b", base, "Janus url base")
	flag.StringVar(&pref, "p", pref, "Janus url path prefix")
	flag.IntVar(&adminPort, "ap", adminPort, "Janus admin endpoint port")
	flag.StringVar(&adminSecret, "s", adminSecret, "Janus admin secret")

	flag.Parse()
	adminEndpoint = fmt.Sprintf("%s:%d", strings.TrimRight(base, "/"), adminPort) + path.Join("/"+pref, "admin")
}

func main() {
	if flag.NArg() == 0 {
		flag.PrintDefaults()
		return
	}
	cmd := flag.Arg(0)
	var cmdFunc func([]string) error
	switch cmd {
	case "videoroom":
		cmdFunc = videoroomCmd
	}
	if cmdFunc == nil {
		fmt.Println("command not found")
	}
	err := cmdFunc(os.Args[indexOf(os.Args, cmd)+1:])
	if err != nil {
		log.Fatal(err)
	}
}

func videoroomCmd(args []string) error {
	flags := flag.NewFlagSet("videoroom", flag.ExitOnError)
	err := flags.Parse(args)
	if err != nil {
		return err
	}
	subcommad := flags.Arg(0)
	if subcommad != "list" {
		return fmt.Errorf("command not found: %s", subcommad)
	}
	tr := &transport.HTTP{
		URL: adminEndpoint,
	}
	acli := jango.Admin{
		Transport:   tr,
		AdminSecret: adminSecret,
	}
	resp := &videoroom.ListResponse{}
	err = acli.PluginRequest(&videoroom.ListRequest{}, resp)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", *resp)

	return nil
}

func indexOf(slice []string, item string) int {
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}
