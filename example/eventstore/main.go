package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-cqrses/cqrses/bus"
	"github.com/go-cqrses/cqrses/esbridge"
	"github.com/go-cqrses/cqrses/eventstore"
	"github.com/go-cqrses/cqrses/eventstore/mysql"
)

func main() {
	ctx := context.Background()

	var es eventstore.EventStore
	es, err := mysql.New(ctx, "root:abcd@tcp(localhost:3306)/events", mysql.DefaultBatchSize)
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err)
	}

	// Clean up any event stream that already exists.
	if err := es.Delete(ctx, "users"); err != nil && err != eventstore.ErrStreamDoesNotExist {
		log.Fatalf("unable to delete users stream: %s", err)
	}

	if err := es.Create(ctx, eventstore.EmptyStreamWithName("users")); err != nil {
		log.Fatalf("unable to create users stream: %s", err)
	}

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

		cmd, _ = changeUserPassword(ctx, "638d863b-3248-4b56-9d0a-e25f62c8cb"+user, "moreSecure")
		if err := cmdBus.Handle(ctx, cmd); err != nil {
			log.Fatalf("unable to change user(%s) to a more secure password: %s", user, err)
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
				if err := events.Next(ctx); err != nil {
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
		case "show events backwards":
			fmt.Println("Showing all recorded events in users stream (backwards):")
			events := es.LoadReverse(ctx, "users", 0, 0, nil)
			for {
				if err := events.Next(ctx); err != nil {
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
				if err := events.Next(ctx); err != nil {
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
