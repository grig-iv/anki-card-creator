package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/alexflint/go-arg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grig-iv/anki-card-creator/creator"
	"github.com/grig-iv/anki-card-creator/ld"
	"golang.org/x/net/proxy"
)

const (
	logPath = "log"
)

var args struct {
	Search string `arg:"-s"`
}

func main() {
	arg.MustParse(&args)

	f, err := tea.LogToFile(logPath, "")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	os.Truncate(logPath, 0)

	proxyClient := newHttpProxyClient()
	ld.HttpClient = proxyClient
	creator.HttpClient = proxyClient

	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func newHttpProxyClient() *http.Client {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
	if err != nil {
		panic(err)
	}

	dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.Dial(network, address)
	}

	transport := &http.Transport{DialContext: dialContext, DisableKeepAlives: true}

	return &http.Client{Transport: transport}
}
