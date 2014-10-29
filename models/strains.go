package models

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/thermokarst/bactdb/router"
)

// A Strain is a subclass of species
type Strain struct {
	Id             int64     `json:"id,omitempty"`
	SpeciesId      int64     `db:"species_id" json:"species_id"`
	StrainName     string    `db:"strain_name" json:"strain_name"`
	StrainType     string    `db:"strain_type" json:"strain_type"`
	Etymology      string    `db:"etymology" json:"etymology"`
	AccessionBanks string    `db:"accession_banks" json:"accession_banks"`
	GenbankEmblDdb string    `db:"genbank_embl_ddb" json:"genbank_eml_ddb"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt      time.Time `db:"deleted_at" json:"deleted_at"`
}

func NewStrain() *Strain {
	return &Strain{StrainName: "Test Strain", StrainType: "Test Type", Etymology: "Test Etymology", AccessionBanks: "Test Accession", GenbankEmblDdb: "Test Genbank"}
}

// StrainService interacts with the strain-related endpoints in bactdb's API
type StrainsService interface {
	// Get a strain
	Get(id int64) (*Strain, error)

	// List all strains
	List(opt *StrainListOptions) ([]*Strain, error)

	// Create a strain record
	Create(strain *Strain) (bool, error)
}

var (
	ErrStrainNotFound = errors.New("strain not found")
)

type strainsService struct {
	client *Client
}

func (s *strainsService) Get(id int64) (*Strain, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.Strain, map[string]string{"Id": strId}, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var strain *Strain
	_, err = s.client.Do(req, &strain)
	if err != nil {
		return nil, err
	}

	return strain, nil
}

func (s *strainsService) Create(strain *Strain) (bool, error) {
	url, err := s.client.url(router.CreateStrain, nil, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("POST", url.String(), strain)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &strain)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusCreated, nil
}

type StrainListOptions struct {
	ListOptions
}

func (s *strainsService) List(opt *StrainListOptions) ([]*Strain, error) {
	url, err := s.client.url(router.Strains, nil, opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	var strains []*Strain
	_, err = s.client.Do(req, &strains)
	if err != nil {
		return nil, err
	}

	return strains, nil
}

type MockStrainsService struct {
	Get_    func(id int64) (*Strain, error)
	List_   func(opt *StrainListOptions) ([]*Strain, error)
	Create_ func(strain *Strain) (bool, error)
}

var _ StrainsService = &MockStrainsService{}

func (s *MockStrainsService) Get(id int64) (*Strain, error) {
	if s.Get_ == nil {
		return nil, nil
	}
	return s.Get_(id)
}

func (s *MockStrainsService) Create(strain *Strain) (bool, error) {
	if s.Create_ == nil {
		return false, nil
	}
	return s.Create_(strain)
}

func (s *MockStrainsService) List(opt *StrainListOptions) ([]*Strain, error) {
	if s.List_ == nil {
		return nil, nil
	}
	return s.List_(opt)
}
