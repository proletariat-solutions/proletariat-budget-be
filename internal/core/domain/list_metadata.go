package domain

type ListMetadata struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// NewListMetadata creates a new ListMetadata instance
func NewListMetadata(total, limit, offset int) *ListMetadata {
	return &ListMetadata{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
}

// HasNextPage returns true if there are more items available
func (m *ListMetadata) HasNextPage() bool {
	return m.Offset+m.Limit < m.Total
}

// HasPreviousPage returns true if there are previous items available
func (m *ListMetadata) HasPreviousPage() bool {
	return m.Offset > 0
}

// GetNextOffset returns the offset for the next page
func (m *ListMetadata) GetNextOffset() int {
	if m.HasNextPage() {
		return m.Offset + m.Limit
	}

	return m.Offset
}

// GetPreviousOffset returns the offset for the previous page
func (m *ListMetadata) GetPreviousOffset() int {
	if m.HasPreviousPage() {
		previousOffset := m.Offset - m.Limit
		if previousOffset < 0 {
			return 0
		}

		return previousOffset
	}

	return m.Offset
}

// GetCurrentPage returns the current page number (1-based)
func (m *ListMetadata) GetCurrentPage() int {
	if m.Limit == 0 {
		return 1
	}

	return (m.Offset / m.Limit) + 1
}

// GetTotalPages returns the total number of pages
func (m *ListMetadata) GetTotalPages() int {
	if m.Limit == 0 {
		return 1
	}

	return (m.Total + m.Limit - 1) / m.Limit
}

// IsValidPagination checks if the pagination parameters are valid
func (m *ListMetadata) IsValidPagination() bool {
	return m.Limit > 0 && m.Offset >= 0 && m.Total >= 0
}
