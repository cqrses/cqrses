# MySQl Event Store

This event store only works with aggregates. We extract aggregate event information from the event and use it to create a unique constratint.

In the future we may introducer "stragagies" that would allow for a custom stream tables.

## Projections

```golang
package main

func main() {
    pm := NewProjectionManager(es)
    sp. err := pm.Create(ctx, "users_dto", []projection.ProjectorOpt{})
    sp.From("users").WhenAny(func (context.Context, messages.Message) error {
        //
    })
    if err := sp.Run(ctx); err != nil {
        log.Fatal(err)
    }
}

```

