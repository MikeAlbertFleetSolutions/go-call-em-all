package callemall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

// List is an abbrevated struct for call-em-all contact list
type List struct {
	ListID   int32  `json:",omitempty"`
	ListType string `json:",omitempty"`
	ListName string
}

// queryListsReturn is the return from 'GET /v1/lists'
type queryListsReturn struct {
	Size     int32
	Items    []List
	Previous string
	Next     string
}

// GetAllLists returns all lists defined
func (client *Client) GetAllLists() ([]List, error) {
	var lists []List

	// handle pagination of results
	url := fmt.Sprintf("%s/v1/lists", client.endpoint)
	for {
		// get batch of lists
		data, err := client.makeRequest("GET", url, nil)
		if err != nil {
			log.Printf("%+v", err)
			return nil, err
		}

		// unmarshal batch of lists
		var ls queryListsReturn
		err = json.Unmarshal(data, &ls)
		if err != nil {
			log.Printf("%+v", err)
			return nil, err
		}

		// append to returned lists
		lists = append(lists, ls.Items...)

		// keep doing til we get them all
		if ls.Next == "" {
			break
		}

		// next pages endpoint
		url = fmt.Sprintf("%s%s", client.endpoint, ls.Next)
	}

	return lists, nil
}

// CreateList creates a single list
func (client *Client) CreateList(name string) (List, error) {
	// create json request body
	l := List{ListName: name}
	b, err := json.Marshal(l)
	if err != nil {
		log.Printf("%+v", err)
		return List{}, err
	}

	// create list, get returned list object
	var data []byte
	data, err = client.makeRequest("POST", fmt.Sprintf("%s/v1/lists", client.endpoint), bytes.NewBuffer(b))
	if err != nil {
		log.Printf("%+v", err)
		return List{}, err
	}

	// unmarshal list
	var list List
	err = json.Unmarshal(data, &list)
	if err != nil {
		log.Printf("%+v", err)
		return List{}, err
	}

	return list, err
}

// DeleteList removes a single list defined
func (client *Client) DeleteList(listID int32) error {
	_, err := client.makeRequest("DELETE", fmt.Sprintf("%s/v1/lists/%d", client.endpoint, listID), nil)
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	return nil
}

// DeleteAllLists removes all lists defined
func (client *Client) DeleteAllLists() error {
	// get list of all lists
	lists, err := client.GetAllLists()
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	// and individually delete
	for _, l := range lists {
		if l.ListType != "System" {
			err = client.DeleteList(l.ListID)
			if err != nil {
				log.Printf("%+v", err)
				return err
			}
		}
	}

	// empty the deleted folder
	_, err = client.makeRequest("DELETE", fmt.Sprintf("%s/v1/lists/deleted/contacts", client.endpoint), nil)
	if err != nil {
		log.Printf("%+v", err)
		return err
	}

	return nil
}
