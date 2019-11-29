// Package ec is an interface for event tracker
package ec

// EventCollector collects the event message.
type EventCollector interface {
	Collect(*Event) error
	CollectByTopic(string, *Event) error
	Close() error
}
