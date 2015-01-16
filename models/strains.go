package models

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/thermokarst/bactdb/router"
)

// A Strain is a subclass of species
type Strain struct {
	Id             int64      `json:"id,omitempty"`
	SpeciesId      int64      `db:"species_id" json:"speciesId"`
	StrainName     string     `db:"strain_name" json:"strainName"`
	StrainType     string     `db:"strain_type" json:"strainType"`
	Etymology      NullString `db:"etymology" json:"etymology"`
	AccessionBanks string     `db:"accession_banks" json:"accessionBanks"`
	GenbankEmblDdb NullString `db:"genbank_embl_ddb" json:"genbankEmblDdb"`
	IsolatedFrom   NullString `db:"isolated_from" json:"isolatedFrom"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt      NullTime   `db:"deleted_at" json:"deletedAt"`
}

type StrainJSON struct {
	Strain *Strain `json:"strain"`
}

type StrainsJSON struct {
	Strains []*Strain `json:"strains"`
}

func (s *Strain) String() string {
	return fmt.Sprintf("%v", *s)
}

func NewStrain() *Strain {
	return &Strain{
		StrainName: "Test Strain",
		StrainType: "Test Type",
		Etymology: NullString{
			sql.NullString{
				String: "Test Etymology",
				Valid:  true,
			},
		},
		AccessionBanks: "Test Accession",
		GenbankEmblDdb: NullString{
			sql.NullString{
				String: "Test Genbank",
				Valid:  true,
			},
		},
		IsolatedFrom: NullString{
			sql.NullString{
				String: "",
				Valid:  false,
			},
		},
	}
}

// StrainService interacts with the strain-related endpoints in bactdb's API
type StrainsService interface {
	// Get a strain
	Get(id int64) (*Strain, error)

	// List all strains
	List(opt *StrainListOptions) ([]*Strain, error)

	// Create a strain record
	Create(strain *Strain) (bool, error)

	// Update an existing strain
	Update(id int64, strain *Strain) (updated bool, err error)

	// Delete an existing strain
	Delete(id int64) (deleted bool, err error)
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

	var strain *StrainJSON
	_, err = s.client.Do(req, &strain)
	if err != nil {
		return nil, err
	}

	return strain.Strain, nil
}

func (s *strainsService) Create(strain *Strain) (bool, error) {
	url, err := s.client.url(router.CreateStrain, nil, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("POST", url.String(), StrainJSON{Strain: strain})
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
	Genus string
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

	var strains *StrainsJSON
	_, err = s.client.Do(req, &strains)
	if err != nil {
		return nil, err
	}

	return strains.Strains, nil
}

func (s *strainsService) Update(id int64, strain *Strain) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.UpdateStrain, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("PUT", url.String(), StrainJSON{Strain: strain})
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req, &strain)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (s *strainsService) Delete(id int64) (bool, error) {
	strId := strconv.FormatInt(id, 10)

	url, err := s.client.url(router.DeleteStrain, map[string]string{"Id": strId}, nil)
	if err != nil {
		return false, err
	}

	req, err := s.client.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return false, err
	}

	var strain *Strain
	resp, err := s.client.Do(req, &strain)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

type MockStrainsService struct {
	Get_    func(id int64) (*Strain, error)
	List_   func(opt *StrainListOptions) ([]*Strain, error)
	Create_ func(strain *Strain) (bool, error)
	Update_ func(id int64, strain *Strain) (bool, error)
	Delete_ func(id int64) (bool, error)
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

func (s *MockStrainsService) Update(id int64, strain *Strain) (bool, error) {
	if s.Update_ == nil {
		return false, nil
	}
	return s.Update_(id, strain)
}

func (s *MockStrainsService) Delete(id int64) (bool, error) {
	if s.Delete_ == nil {
		return false, nil
	}
	return s.Delete_(id)
}
