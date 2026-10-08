package repository

import (
	"contact-manager/model"
)

type ContactRepository struct {
	contacts []model.Contact
	nextID   int
}

func NewContactRepository() *ContactRepository {
	return &ContactRepository{
		contacts: make([]model.Contact, 0),
		nextID: 1,
	}
}

func (r *ContactRepository) Add(contact model.Contact) model.Contact {
	contact.ID = r.nextID
	r.nextID++

	r.contacts = append(r.contacts, contact)

	return contact
}

func (r *ContactRepository) GetAll() []model.Contact {
	return r.contacts
}

func (r *ContactRepository) GetByID(id int) (*model.Contact, bool) {
	for i := range r.contacts {
		if r.contacts[i].ID == id {
			return &r.contacts[i], true
		}
	}
	return nil, false
}

func (r *ContactRepository) UpdateById(id int, updatedContact model.Contact) bool {
	for i :=range r.contacts {
		if r.contacts[i].ID==id {
			r.contacts[i] = updatedContact
			return true
		}
	}
	return false
}

func (r *ContactRepository) Delete(id int) bool {
	for i := range r.contacts {
		if r.contacts[i].ID == id {
			r.contacts = append(r.contacts[:i], r.contacts[i+1:]...)
			return true
		}
	}
	return false
}
