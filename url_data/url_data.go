package url_data

import (
	"encoding/json"
	"fmt"
	"time"
)

// basic url record struct
// both json and bson tags are required:
// json handles http api part,
// bson handles db part
// omitempty is required for db filters
type URLData struct {
	ID        string `json:"_id,omitempty" bson:"_id,omitempty"`
	URL       string `json:"url" bson:"url,omitempty"` // json.url cannot be empty
	ShortCode string `json:"shortCode,omitempty" bson:"shortCode,omitempty"`
	// custom-marshaled propeties
	CreatedAt   time.Time `json:"-" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"-" bson:"updatedAt,omitempty"`
	AccessCount int       `json:"-" bson:"accessCount,omitempty"`
	// control properties
	include_access_count_in_json bool `json:"-" bson:"-"`
}

// alias to avoid recursion during marshal/unmarshal
type urlDataAlias URLData

// auxiliary type for marshal/unmarshal (doesn't have MarshalJSON/UnmarshalJSON methods)
type urlDataAux struct {
	*urlDataAlias        // embed all fields from URLData
	CreatedAt     string `json:"createdAt,omitempty"`   // shadowed
	UpdatedAt     string `json:"updatedAt,omitempty"`   // shadowed
	AccessCount   *int   `json:"accessCount,omitempty"` // use pointer, json handles them gracefully
}

// controls whether to include access count in json or not
func (u *URLData) IncludeAccessCountInJSON(include bool) {
	u.include_access_count_in_json = include
}

// json marshaler (convert to []byte)
func (u *URLData) MarshalJSON() ([]byte, error) {

	ac_val := &u.AccessCount
	if !u.include_access_count_in_json {
		ac_val = nil // exclude accessCount
	}
	aux := &urlDataAux{
		urlDataAlias: (*urlDataAlias)(u), // convert URLData to Alias type
		CreatedAt:    u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    u.UpdatedAt.Format(time.RFC3339),
		AccessCount:  ac_val,
	}
	return json.Marshal(aux)
}

// json unmarshaler (convert from []byte)
func (u *URLData) UnmarshalJSON(data []byte) error {

	aux := &urlDataAux{
		urlDataAlias: (*urlDataAlias)(u),
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	// parse custom date to time.Time
	u.CreatedAt, err = time.Parse(time.RFC3339, aux.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed parsing createdAt: %v", err)
	}
	u.UpdatedAt, err = time.Parse(time.RFC3339, aux.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed parsing updatedAt: %v", err)
	}

	return nil
}

// stringer for URLData
func (u URLData) String() string {
	res, err := json.Marshal(&u)
	if err != nil {
		return ""
	}
	return string(res)
}
