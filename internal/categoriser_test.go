package internal

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"
)

func TestCategoriserProcessing(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ctg, err := NewCategoriser(log.Default(), []CategoriserConfig{
		{SourceRegexp: "/v1/auth/(sign-in|sign-out|login).*", Target: "/v1/auth/$1"},
		{SourceRegexp: "/v2/([a-zA-Z0-9]+)/([a-zA-Z0-9]+).*", Target: "/v2/$1/$2"},
	})

	t.Log("categoriser initialized")

	if err != nil {
		t.Fatal("failed to create categoriser service")
	}

	t.Log("starting processing")
	inputC := make(chan SyslogMessage)
	outputC := make(chan NginxEvent)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		t.Log("started subprocess")
		defer wg.Done()

		ctg.Run(ctx, inputC, outputC)
		t.Log("categoriser exited")
	}()

	testMap := map[string]string{
		"/v1/auth/sign-in/auth":     "/v1/auth/sign-in",
		"/v2/testing/test/some/url": "/v2/testing/test",
		"/v2/some":                  "other",
	}

	for k, v := range testMap {
		inputC <- SyslogMessage{
			Message: `{"request":"` + k + `"}`,
		}

		t.Log("message sent, waiting for reply")

		select {
		case parsed := <-outputC:
			if parsed.Category != v {
				t.Logf("incorrect category after analyze: %s != %s", k, v)
			} else {
				t.Logf("remapped category %s => %s", k, parsed.Category)
			}
		case <-ctx.Done():
			t.Fatal("failed to handle during timeout")
		}
	}
}
