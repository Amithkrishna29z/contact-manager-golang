package service

import (
	"contact-manager/model"
	"contact-manager/repository"
	"errors"
)

type ContactService struct {
	repository *repository.ContactRepository
}

func NewContactService(
	repository *repository.ContactRepository,
) *ContactService {

	return &ContactService {
		repository: repository,
	}
}

func (s *ContactService) AddContact(
	name string,
	phone string,
) (model.Contact, error) {

	if name == "" {
		return model.Contact{}, errors.New("Name cannot be empty")
	}
	if phone == "" {
		return model.Contact{}, errors.New("Phone cannot be empty")
	}

	contact := model.Contact {
		Name: name,
		Phone: phone,
	}
	return s.repository.Add(contact),nil
}

func(s *ContactService) GetContacts() []model.Contact {
	return s.repository.GetAll()
}

func (s *ContactService) GetContact(id int) (model.Contact, error)  {
	contact, found:=s.repository.GetByID(id)

	if !found {
		return model.Contact{}, errors.New("Contact not found")
	}
	return *contact, nil
}

func (s *ContactService) DeleteContact(id int) error {
	deleted := s.repository.Delete(id)

	if !deleted {
		return errors.New("Contact not found")
	}
	return nil
}