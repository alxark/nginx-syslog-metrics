package internal

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"regexp"
)

type compiledRegexp struct {
	source *regexp.Regexp
	target string
}

type Categoriser struct {
	cfg []compiledRegexp
	log *log.Logger
}

func NewCategoriser(log *log.Logger, config []CategoriserConfig) (cr *Categoriser, err error) {
	cr = &Categoriser{
		log: log,
	}

	for _, v := range config {
		if v.SourceRegexp == "" || v.Target == "" {
			return cr, errors.New("failed to compiled regexp, it's empty")
		}

		cr.cfg = append(cr.cfg, compiledRegexp{
			source: regexp.MustCompile(v.SourceRegexp),
			target: v.Target,
		})

		cr.log.Printf("compiled category: %s", v.SourceRegexp)
	}

	cr.log.Printf("total categories compiled: %d", len(cr.cfg))

	return cr, nil
}

func (cr *Categoriser) Run(ctx context.Context, input chan SyslogMessage, output chan NginxEvent) error {
	for ctx.Err() == nil {
		select {
		case msg := <-input:
			var ne NginxEvent

			if err := json.Unmarshal([]byte(msg.Message), &ne); err != nil {
				continue
			}

			ne.Frontend = msg.Host

			for _, v := range cr.cfg {
				if !v.source.MatchString(ne.Request) {
					continue
				}

				ne.Category = v.source.ReplaceAllString(ne.Request, v.target)
				break
			}

			if ne.Category == "" {
				ne.Category = "other"
				cr.log.Printf(`failed to detect URL category: %s`, ne.Request)
			}

			output <- ne
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
