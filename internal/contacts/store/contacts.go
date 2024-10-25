package store

import (
	"github.com/gavink97/templ-campaigner/internal/contacts"
	"github.com/gavink97/templ-campaigner/internal/handlers"
	"gorm.io/gorm"
)

type ContactStore struct {
	db *gorm.DB
}

type NewContactStoreParams struct {
	DB *gorm.DB
}

func NewContactStore(params NewContactStoreParams) *ContactStore {
	return &ContactStore{
		db: params.DB,
	}
}

func (s *ContactStore) CreateContact(fname, lname, email string, subscribed, unsubscribed bool) error {
	fname = handlers.MakeTitle(fname)
	lname = handlers.MakeTitle(lname)

	return s.db.Create(&contacts.Contact{
		FName:        fname,
		LName:        lname,
		EmailAddress: email,
		Subscribed:   subscribed,
		Unsubscribed: unsubscribed,
	}).Error
}

func (s *ContactStore) GetContact(email string) (*contacts.Contact, error) {
	var user contacts.Contact
	err := s.db.Where("email_address = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, err
}

func (s *ContactStore) GetSubscribersList() (*[]contacts.Contact, error) {
	var list []contacts.Contact
	err := s.db.Where("subscribed = true").Find(&list).Error

	if err != nil {
		return nil, err
	}
	return &list, err
}

func (s *ContactStore) SearchContacts(str string) (*[]contacts.Contact, error) {
	var users []contacts.Contact
	wildstr := str + "%"
	err := s.db.Order("email_address, f_name, l_name").Where(`email_address LIKE
        ?`, wildstr).Or("f_name LIKE ?", wildstr).Or("l_name LIKE ?", wildstr).Find(&users).Error

	if err != nil {
		return nil, err
	}

	return &users, err
}
