package messages

type metaKey string

const (
	// MetaCausationID should be set on a context where a command
	// is called, this will allow us to pull the command id out and
	// set it as the causation ID.
	// Read more: https://blog.arkency.com/correlation-id-and-causation-id-in-evented-systems/
	MetaCausationID metaKey = "causation_id"

	// MetaCorrelationID if another event caused this new event we should
	// set that as the correlation id on the event's metadata.
	// Read more: https://blog.arkency.com/correlation-id-and-causation-id-in-evented-systems/
	MetaCorrelationID metaKey = "correlation_id"

	// MetaAggregateID is the identifier of an aggregate, used when using event sourcing.
	MetaAggregateID metaKey = "aggregate_id"

	// MetaAggregateVersion should be set on an event for the version it is targeting, for
	// example if you've loaded an aggregate at version 6 the next version should be 7.
	// Versions start from 1!
	MetaAggregateVersion metaKey = "aggregate_version"
)
