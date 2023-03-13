package internal

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Application struct {
	SyslogPort      int
	QueueSize       int
	HttpPort        int
	Prefix          string
	CategoriserFile string
	log             *log.Logger
}

func NewApplication(logger *log.Logger) *Application {
	a := &Application{}
	a.log = logger

	return a
}

func (a *Application) ConfigureFlags() error {
	flag.IntVar(&a.SyslogPort, "syslogPort", 9090, "UDP port used to receive syslog messages")
	flag.IntVar(&a.QueueSize, "queueSize", 8192, "size of internal queue for received messages")
	flag.IntVar(&a.HttpPort, "httpPort", 8080, "HTTP port used to read metrics")
	flag.StringVar(&a.Prefix, "prefix", "nsm", "prometheus variable prefix")
	flag.StringVar(&a.CategoriserFile, "categoriserFile", "", "categoriser file, should be a JSON file")

	flag.Parse()

	return nil
}

func (a *Application) Run() error {
	ctx := context.Background()

	a.log.Printf("starting service, syslogPort: %d, httpPort: %d", a.SyslogPort, a.HttpPort)

	var catCfg []CategoriserConfig
	if a.CategoriserFile != "" {
		f, err := os.Open(a.CategoriserFile)
		if err != nil {
			return err
		}

		catCfgBody, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		err = json.Unmarshal(catCfgBody, &catCfg)
		if err != nil {
			return err
		}
	}

	rcvr, err := NewReceiver(a.log, a.SyslogPort, a.QueueSize, a.Prefix)
	if err != nil {
		return err
	}

	httpServer, err := NewHttpServer(a.log, a.HttpPort)
	if err != nil {
		return err
	}

	categoriser, err := NewCategoriser(a.log, catCfg)
	if err != nil {
		return err
	}

	counter, err := NewCounter(a.log, a.Prefix)
	if err != nil {
		return err
	}

	syslogMessages := make(chan SyslogMessage, a.QueueSize)
	nginxMessages := make(chan NginxEvent, a.QueueSize)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		a.log.Print("started syslog thread")
		defer wg.Done()
		if err := rcvr.Run(syslogMessages); err != nil {
			a.log.Fatalf("receiver is failed, got %s", err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		a.log.Print("started HTTP-server thread")
		defer wg.Done()
		if err := httpServer.Run(); err != nil {
			a.log.Fatalf("failed to run HTTP server, got %s", err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		a.log.Print("starting categoriser thread")
		defer wg.Done()

		if err := categoriser.Run(ctx, syslogMessages, nginxMessages); err != nil {
			a.log.Fatalf("failed to handle categoriser, got %s", err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		a.log.Print("starting counter for metrics")
		defer wg.Done()

		if err := counter.Run(ctx, nginxMessages); err != nil {
			a.log.Fatalf("failed to count events, got %s", err.Error())
		}
	}()

	a.log.Print("waiting for subroutines to complete")
	wg.Wait()

	return nil
}
