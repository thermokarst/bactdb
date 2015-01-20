package models

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/thermokarst/bactdb/router"
)

// A Genus is a high-level classifier in bactdb.
type GenusBase struct {
	Id        int64     `json:"id,omitempty"`
	GenusName string    `db:"genus_name" json:"genusName"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt NullTime  `db:"deleted_at" json:"deletedAt"`
}

type Genus struct {
	*GenusBase
	Species NullSliceInt64 `db:"species" json:"species"`
}

type GenusJSON struct {
	Genus *Genus `json:"genus"`
}

type GeneraJSON struct {
	Genera []*Genus `json:"genera"`
}

func (m *Genus) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *GenusBase) String() string {
	return fmt.Sprintf("%v", *m)
}

func NewGenus() *Genus {
	return &Genus{&GenusBase{GenusName: "Test Genus"}, make([]int64, 0)}
}

// GeneraService interacts with the genus-related endpoints in bactdb's API.
type GeneraService interface {
	// Get a genus.
	Get(id int64) (*Genus, error)

	// List all genera.
	List(opt *GenusListOptions) ([]*Genus, error)

	// Create a new genus. The newly created genus's ID is written to genus.Id
	Create(genus *Genus) (created bool, err error)

	// Update an existing genus.
	Update(id int64, genus *Genus) (updated bool, err error)

	// Delete an existing genus.
	Delete(id int64) (deleted bool, err error)
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

	var genus *GenusJSON
	_, err = s.client.Do(req, &genus)
	if err != nil {
		return nil, err
	}

	return genus.Genus, nil
}

func (s *generaService) Create(genus *Genus) (bool, error) {
	url, err := s.client.url(router.CreateGenus, nil, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("POST", url.String(), GenusJSON{Genus: genus})
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

	var genera *GeneraJSON
	_, err = s.client.Do(req, &genera)
	if err != nil {
		return nil, err
	}

	return genera.Genera, nil
}

func (s *generaService) Update(id int64, genus *Genus) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.UpdateGenus, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("PUT", url.String(), GenusJSON{Genus: genus})
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &genus)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (s *generaService) Delete(id int64) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.DeleteGenus, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return false, err
	}

	var genus *Genus
	resp, err := s.client.Do(req, &genus)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

type MockGeneraService struct {
	Get_    func(id int64) (*Genus, error)
	List_   func(opt *GenusListOptions) ([]*Genus, error)
	Create_ func(genus *Genus) (bool, error)
	Update_ func(id int64, genus *Genus) (bool, error)
	Delete_ func(id int64) (bool, error)
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

func (s *MockGeneraService) Update(id int64, genus *Genus) (bool, error) {
	if s.Update_ == nil {
		return false, nil
	}
	return s.Update_(id, genus)
}

func (s *MockGeneraService) Delete(id int64) (bool, error) {
	if s.Delete_ == nil {
		return false, nil
	}
	return s.Delete_(id)
}
