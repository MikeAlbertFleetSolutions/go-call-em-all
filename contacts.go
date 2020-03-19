package callemall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

// Contact is an abbrevated struct for call-em-all contacts
type Contact struct {
	PersonID     int32 `json:",omitempty"`
	FirstName    string
	LastName     string
	PrimaryPhone string
	Lists        []List `json:",omitempty"`
}

// CreateContact creates a single contact in the account
func (client *Client) CreateContact(firstname, lastname, phone string, lists []List) (Contact, error) {
	// create json request body
	c := Contact{
		FirstName:    firstname,
		LastName:     lastname,
		PrimaryPhone: phone,
		Lists:        lists,
	}
	b, err := json.Marshal(c)
	if err != nil {
		log.Printf("%+v", err)
		return Contact{}, err
	}

	// create contact, get returned contact object
	var data []byte
	data, err = client.makeRequest("POST", fmt.Sprintf("%s/v1/contacts", client.endpoint), bytes.NewBuffer(b))
	if err != nil {
		log.Printf("%+v", err)
		return Contact{}, err
	}

	// unmarshal contact
	var contact Contact
	err = json.Unmarshal(data, &contact)
	if err != nil {
		log.Printf("%+v", err)
		return Contact{}, err
	}

	return contact, nil
}

// DeleteContact removes a single contact defined in the account
func (client *Client) DeleteContact(personID int32) error {
	_, err := client.makeRequest("DELETE", fmt.Sprintf("%s/v1/contacts/%d", client.endpoint, personID), nil)
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	return nil
}

// DeleteAllContacts removes all contacts defined in the account
func (client *Client) DeleteAllContacts() error {
	_, err := client.makeRequest("DELETE", fmt.Sprintf("%s/v1/contacts", client.endpoint), nil)
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	return nil
}
