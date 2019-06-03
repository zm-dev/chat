package memory_store

import (
	"github.com/zm-dev/chat/model"
	"sync"
)

type Ticket struct {
	tickets sync.Map
}

func (t *Ticket) TicketIsNotExistErr(err error) bool {
	return model.TicketIsNotExistErr(err)
}

func (t *Ticket) TicketLoad(id string) (*model.Ticket, error) {

	if ticket, ok := t.tickets.Load(id); ok {
		return ticket.(*model.Ticket), nil
	}
	return nil, model.ErrTicketNotExist
}

func (t *Ticket) TicketCreate(ticket *model.Ticket) error {

	if _, ok := t.tickets.Load(ticket.Id); ok {
		return model.ErrTicketExisted
	}
	t.tickets.Store(ticket.Id, ticket)
	return nil

}

func (t *Ticket) TicketDelete(id string) error {
	if _, ok := t.tickets.Load(id); !ok {
		return model.ErrTicketNotExist
	}
	t.tickets.Delete(id)
	return nil
}

func NewMemoryTicket() model.TicketStore {
	return &Ticket{
		sync.Map{},
	}
}
