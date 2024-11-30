package backend

import (
	"fmt"
	"url-shortener/db_interface"
)

// mock db interface

type dbCollectionMock struct {
	data   []URLData
	id_cnt int
}

func (collection *dbCollectionMock) InsertOne(doc any) (id string, err error) {
	t, ok := doc.(URLData)
	if ok {
		t.ID = fmt.Sprintf("%d", collection.id_cnt)
		collection.data = append(collection.data, t)
		collection.id_cnt++
		return t.ID, nil
	}
	return "", fmt.Errorf("invalid doc type %T", t)
}

func (collection *dbCollectionMock) FindOne(filter any, result any) error {
	f, ok := filter.(URLData)
	if !ok {
		return fmt.Errorf("invalid filter type %T", f)
	}
	r, ok := result.(*URLData)
	if !ok {
		return fmt.Errorf("invalid result type %T", r)
	}
	for _, data := range collection.data {
		if f.URL == data.URL || f.ShortCode == data.ShortCode {
			*r = data
			return nil
		}
	}
	return db_interface.ErrNoDocuments
}

// update doc
func (collection *dbCollectionMock) UpdateOne(filter any, update_with any) error {
	f, ok := filter.(URLData)
	if !ok {
		return fmt.Errorf("invalid filter type %T", f)
	}
	r, ok := update_with.(*URLData)
	if !ok {
		return fmt.Errorf("invalid result type %T", r)
	}

	update := func(first *URLData, second *URLData) {
		if second.URL != "" {
			first.URL = second.URL
		}
		if second.ShortCode != "" {
			first.ShortCode = second.ShortCode
		}
		if second.ID != "" {
			first.ID = second.ID
		}
		if !second.CreatedAt.IsZero() {
			first.CreatedAt = second.CreatedAt
		}
		if !second.UpdatedAt.IsZero() {
			first.UpdatedAt = second.UpdatedAt
		}
		if second.AccessCount != 0 {
			first.AccessCount = second.AccessCount
		}
	}

	for i := range collection.data {
		data := &collection.data[i]
		if f.URL == data.URL || f.ShortCode == data.ShortCode {
			update(data, r)
			update(r, data)
			return nil
		}
	}
	return db_interface.ErrNoDocuments
}

// delete doc
func (collection *dbCollectionMock) DeleteOne(filter any) error {
	f, ok := filter.(URLData)
	if !ok {
		return fmt.Errorf("invalid filter type %T", f)
	}
	for i, data := range collection.data {
		if f.URL == data.URL || f.ShortCode == data.ShortCode {
			temp := collection.data[i+1:]                      // save everything after i
			collection.data = collection.data[:i]              // truncate until i
			collection.data = append(collection.data, temp...) // concatenate
			return nil
		}
	}
	return db_interface.ErrNoDocuments
}

// find some records
func (collection *dbCollectionMock) FindSome(limit int, result any) error {
	f, ok := result.(*[]URLData)
	if !ok {
		return fmt.Errorf("invalid result type %T", f)
	}
	lim := min(limit, cap(collection.data))
	*f = collection.data[:lim]
	return nil
}
