package models

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/thermokarst/bactdb/router"
)

// A Genus is a high-level classifier in bactdb.
type Genus struct {
	Id        int64     `json:"id,omitempty"`
	GenusName string    `json:"genus_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// GeneraService interacts with the genus-related endpoints in bactdb's API.
type GeneraService interface {
	// Get a genus.
	Get(id int64) (*Genus, error)

	// List all genera.
	List(opt *GenusListOptions) ([]*Genus, error)

	// Create a new genus. The newly created genus's ID is written to genus.Id
	Create(genus *Genus) (created bool, err error)
}

var (
	ErrGenusNotFound = errors.New("genus not found")
)

type generaService struct {
	client *Client
}

func (s *generaService) Get(id int64) (*Genus, error) {
	// Pass in key value pairs as strings, so that the gorilla mux URL
	// generation is happy.
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.Genus, map[string]string{"Id": strId}, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var genus *Genus
	_, err = s.client.Do(req, &genus)
	if err != nil {
		return nil, err
	}

	return genus, nil
}

func (s *generaService) Create(genus *Genus) (bool, error) {
	url, err := s.client.url(router.CreateGenus, nil, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("POST", url.String(), genus)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &genus)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusCreated, nil
}

type GenusListOptions struct {
	ListOptions
}

func (s *generaService) List(opt *GenusListOptions) ([]*Genus, error) {
	url, err := s.client.url(router.Genera, nil, opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var genera []*Genus
	_, err = s.client.Do(req, &genera)
	if err != nil {
		return nil, err
	}

	return genera, nil
}

type MockGeneraService struct {
	Get_    func(id int64) (*Genus, error)
	List_   func(opt *GenusListOptions) ([]*Genus, error)
	Create_ func(post *Genus) (bool, error)
}

var _ GeneraService = &MockGeneraService{}

func (s *MockGeneraService) Get(id int64) (*Genus, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}

func (s *MockGeneraService) Create(genus *Genus) (bool, error) {
	if s.Create_ == nil {
		return false, nil
	}
	return s.Create_(genus)
}

func (s *MockGeneraService) List(opt *GenusListOptions) ([]*Genus, error) {
	if s.List_ == nil {
		return nil, nil
	}
	return s.List_(opt)
}
