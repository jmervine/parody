package main

import (
	"bytes"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var copyClient, mainClient *client
var thisHost, mainHost, mainName, copyHost, copyName string

func init() {
	app := cli.NewApp()
	app.Name = "parody"
	app.Version = "1.0.0"
	app.Author = "Joshua Mervine"
	app.Email = "joshua@mervine.net"
	app.Usage = "simple http proxy for copying posts and puts to two locations"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen,l",
			Usage: "listener address",
			Value: "localhost:8888",
		},
		cli.StringFlag{
			Name:  "main,m",
			Usage: "main upstream location [host:port]",
		},
		cli.StringFlag{
			Name:  "copy,c",
			Usage: "copy upstream location [host:port]",
		},
		cli.StringFlag{
			Name:  "main-name",
			Usage: "main upstream name in logger",
			Value: "main",
		},
		cli.StringFlag{
			Name:  "copy-name",
			Usage: "copy upstream name in logger",
			Value: "copy",
		},
	}

	app.Action = func(c *cli.Context) {
		thisHost = c.String("listen")
		mainHost = c.String("main")
		copyHost = c.String("copy")
		mainName = c.String("main-name")
		copyName = c.String("copy-name")

		if mainHost == "" || copyHost == "" {
			cli.ShowAppHelp(c)
		}
	}

	app.Run(os.Args)
}

func main() {
	if mainHost == "" || copyHost == "" {
		os.Exit(1)
	}

	copyClient = &client{client: &http.Client{}, host: copyHost}
	mainClient = &client{client: &http.Client{}, host: mainHost}

	proxy := http.HandlerFunc(handler)
	log.Println("Listening on", thisHost)
	log.Fatal(http.ListenAndServe(thisHost, proxy))
}

func handler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	req.RequestURI = ""

	body, _ := ioutil.ReadAll(req.Body)

	logger := func(name string, res *http.Response) {
		log.Printf("[%s] %s - %s - %s\n", name, res.Status, req.URL.Path, time.Since(start))
	}

	if req.Method == "POST" || req.Method == "PUT" {
		go func() {
			res, _ := copyClient.do(*req, body)
			logger(copyName, res)
		}()
	}

	res, err := mainClient.do(*req, body)

	defer func() {
		res.Body.Close()
		logger(mainName, res)
	}()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res.Write(w)
}

// client
type client struct {
	client *http.Client
	host   string
}

func (c client) do(req http.Request, body []byte) (*http.Response, error) {
	var err error

	req.URL.Scheme = "http"
	req.URL.Host = c.host
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	req.Header.Set("Connection", "keepalive")
	req.Header.Set("Content-Length", string(len(body)+1))

	res, err := c.client.Do(&req)
	if err != nil {
		log.Printf("ERROR (%s): %s\n", c.host, err)
		return res, err
	}
	return res, err
}
