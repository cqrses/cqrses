package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"gopkg.in/cqrses/bus"
	"gopkg.in/cqrses/esbridge"
	"gopkg.in/cqrses/eventstore"
	"gopkg.in/cqrses/eventstore/inmem"
)

func main() {
	ctx := context.Background()

	var es eventstore.EventStore = inmem.New()
	es.Create(ctx, eventstore.EmptyStreamWithName("users"))

	eventBus := bus.NewEventBus()

	users := userCollection{}
	proj := &userProjector{
		all: users,
	}
	eventBus.Register(bus.MatchAny(), proj.Handle)
	es = eventBus.WrapStore(es)

	cmdBus := bus.NewCommandBus()
	cmdBus.PushMiddleware(esbridge.AttachEventStoreToBus(es))

	registerCommandBusHandlers(cmdBus)

	for _, user := range []string{"a1", "b2", "c3", "d4", "e5"} {
		cmd, _ := createUserWith(ctx, "638d863b-3248-4b56-9d0a-e25f62c8cb"+user, user+"@testing.com", "changeme")

		if err := cmdBus.Handle(ctx, cmd); err != nil {
			log.Fatalf("unable to create user(%s): %s", user, err)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("#start")
		fmt.Print("action: ")
		scanner.Scan()

		switch scanner.Text() {
		case "show users":
			fmt.Println("Showing all users:")
			for _, udto := range users {
				fmt.Printf(
					" - user: (id: %s) (email address: %s)\n",
					udto.id,
					udto.emailAddress,
				)
			}
		case "show events":
			fmt.Println("Showing all recorded events in users stream:")
			events := es.Load(ctx, "users", 0, 0, nil)
			for {
				if err := events.Next(); err != nil {
					if err == eventstore.EOF {
						break
					}

					log.Fatalf("error: %s", err)
				}

				e := events.Current()
				fmt.Printf(
					" - event: %s (id: %s)\n",
					e.MessageName(),
					e.MessageID(),
				)
			}
		case "show raw events":
			fmt.Println("Showing raw recorded events in users stream:")
			events := es.Load(ctx, "users", 0, 0, nil)
			for {
				if err := events.Next(); err != nil {
					if err == eventstore.EOF {
						break
					}

					log.Fatalf("error: %s", err)
				}

				fmt.Printf(" - event: %+v\n", events.Current())
			}
		case "quit":
			fmt.Println("Bye!")
			os.Exit(0)
		default:
			fmt.Println("Command not recognised.")
		}
		fmt.Print("#end\n\n")
	}
}
