package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"gopkg.in/cqrses/eventstore"
	"gopkg.in/cqrses/eventstore/inmem"
	"gopkg.in/cqrses/messages"
)

func main() {
	ctx := context.Background()
	es := inmem.New()

	es.Create(ctx, eventstore.EmptyStreamWithName("test"))

	e1 := messages.NewEvent("hello-world", "EchoText", map[string]interface{}{}, map[string]interface{}{}, 0, time.Now())
	es.AppendTo(ctx, "test", []*messages.Event{e1})

	events := es.Load(ctx, "test", 0, 1, nil)
	for {
		if err := events.Next(); err != nil {
			if err == eventstore.EOF {
				break
			}

			log.Fatalf("error: %s", err)
		}

		e := events.Current()
		fmt.Printf(
			"Received Event: %s (id: %s)",
			e.MessageName(),
			e.MessageID(),
		)
	}
}
