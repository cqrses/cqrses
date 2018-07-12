package aggregate

import (
	"gopkg.in/cqrses/messages"
)

type (
	BaseAggregate struct {
		id      string
		history *History
		apply   func(*messages.Event) error
	}
)

func (a *BaseAggregate) Identifier() string {
	return a.id
}

func (a *BaseAggregate) WithHistory(*History) {
	return
}

func (a *BaseAggregate) History() *History {
	return a.history
}

func (a *BaseAggregate) Apply(*messages.Event) error {

}
