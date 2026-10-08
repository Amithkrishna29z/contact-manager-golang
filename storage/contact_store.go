package storage

import (
	"contact-manager/model"
)

type ContactStore interface {
	Save(contact model.Contact) error
	Load() ([]model.Contact, error)
}
