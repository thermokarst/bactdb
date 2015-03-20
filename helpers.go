package main

// ListOptions specifies general pagination options for fetching a list of results
type ListOptions struct {
	PerPage int   `url:",omitempty" json:",omitempty"`
	Page    int   `url:",omitempty" json:",omitempty"`
	Ids     []int `url:",omitempty" json:",omitempty" schema:"ids[]"`
}

func (o ListOptions) PageOrDefault() int {
	if o.Page <= 0 {
		return 1
	}
	return o.Page
}

func (o ListOptions) Offset() int {
	return (o.PageOrDefault() - 1) * o.PerPageOrDefault()
}

func (o ListOptions) PerPageOrDefault() int {
	if o.PerPage <= 0 {
		return DefaultPerPage
	}
	return o.PerPage
}

// DefaultPerPage is the default number of items to return in a paginated result set
const DefaultPerPage = 10
