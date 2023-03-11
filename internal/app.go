package internal

import (
	"flag"
	"log"
	"sync"
)

type Application struct {
	SyslogPort      int
	SyslogQueueSize int
	HttpPort        int
	Prefix          string
	log             *log.Logger
}

func NewApplication(logger *log.Logger) *Application {
	a := &Application{}
	a.log = logger

	return a
}

func (a *Application) ConfigureFlags() error {
	flag.IntVar(&a.SyslogPort, "syslogPort", 9090, "UDP port used to receive syslog messages")
	flag.IntVar(&a.SyslogQueueSize, "syslogQueueSize", 8192, "size of internal queue for received messages")
	flag.IntVar(&a.HttpPort, "httpPort", 8080, "HTTP port used to read metrics")
	flag.StringVar(&a.Prefix, "prefix", "nsm", "prometheus variable prefix")

	flag.Parse()

	return nil
}

func (a *Application) Run() error {
	a.log.Printf("starting service, syslogPort: %d, httpPort: %d", a.SyslogPort, a.HttpPort)

	incomming, err := NewReceiver(a.log, a.SyslogPort, a.SyslogQueueSize, a.Prefix)
	if err != nil {
		return err
	}

	httpServer, err := NewHttpServer(a.log, a.HttpPort)
	if err != nil {
		return err
	}

	outputChannel := make(chan SyslogMessage, a.SyslogQueueSize)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		a.log.Print("started syslog thread")
		defer wg.Done()
		incomming.Run(outputChannel)
	}()

	wg.Add(1)
	go func() {
		a.log.Print("started HTTP-server thread")
		defer wg.Done()
		httpServer.Run()
	}()

	for v := range outputChannel {
		a.log.Print(v.Message)
	}

	a.log.Print("waiting for subroutines to complete")
	wg.Wait()

	return nil
}
