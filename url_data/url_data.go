package url_data

import (
	"encoding/json"
	"fmt"
	"time"
)

type URLData struct {
	ID          string    `json:"_id,omitempty" bson:"_id,omitempty"`
	URL         string    `json:"url" bson:"url,omitempty"`
	ShortCode   string    `json:"shortCode" bson:"shortCode,omitempty"`
	CreatedAt   time.Time `json:"-" bson:"createdAt,omitempty"` // do not marshal json automatically
	UpdatedAt   time.Time `json:"-" bson:"updatedAt,omitempty"` // do not marshal json automatically
	AccessCount int       `json:"accessCount" bson:"accessCount,omitempty"`
	// control properties
	include_access_count_in_json bool `json:"-" bson:"-"`
}

// controls whether to include access count in json or not
func (u *URLData) IncludeAccessCountInJSON(include bool) {
	u.include_access_count_in_json = include
}

// json marshaler (convert to []byte)
func (u *URLData) MarshalJSON() ([]byte, error) {
	type Alias URLData // alias to avoid recursion
	type Aux struct {
		*Alias             // embed all fields from URLData
		CreatedAt   string `json:"createdAt,omitempty"` // these are shadowed
		UpdatedAt   string `json:"updatedAt,omitempty"`
		AccessCount *int   `json:"accessCount,omitempty"` // use pointer, json handles them gracefully
	}
	ac_val := &u.AccessCount
	if !u.include_access_count_in_json {
		ac_val = nil // exclude accessCount
	}
	aux := &Aux{
		Alias:       (*Alias)(u), // convert URLData to Alias type (removes MarshalJSON as it's not defined for Aux)
		CreatedAt:   u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   u.UpdatedAt.Format(time.RFC3339),
		AccessCount: ac_val,
	}
	return json.Marshal(aux)
}

// json unmarshaler (convert from []byte)
func (u *URLData) UnmarshalJSON(data []byte) error {
	type Alias URLData // alias to avoid recursion
	type Aux struct {
		*Alias
		CreatedAt string `json:"createdAt,omitempty"`
		UpdatedAt string `json:"updatedAt,omitempty"`
	}
	aux := &Aux{
		Alias: (*Alias)(u),
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

func (u URLData) String() string {
	res, err := json.Marshal(&u)
	if err != nil {
		return ""
	}
	return string(res)
}
