package elasticsearch

import (
	"context"
	"github.com/olivere/elastic/v7"
)

type SearchQuery interface {
	Query() elastic.Query
}

type SearchSorter interface {
	Sorter() elastic.Sorter
}

type SearchCursor interface {
	Cursor() []interface{}
}

type Searcher interface {
	IndexName() string
	Size() int
	Query() elastic.Query
	Sorter() []elastic.Sorter
	Cursor() []interface{}
}

type Page struct {
	From int
	Size int
	Desc bool
}

func (client *Client) Search(searcher Searcher, Response interface{}) error {

	// var hits [][]byte

	for {

		resp, err := client.client.Search().
			Index(searcher.IndexName()).
			Size(searcher.Size()).
			Query(searcher.Query()).
			SortBy(searcher.Sorter()...).
			SearchAfter(searcher.Cursor()...).
			Do(context.Background())

		if err != nil {
			return err
		}

		if resp.TotalHits() <= 0 {
			break
		}

		if len(resp.Hits.Hits) == 0 {
			break
		}

	}

	return nil
}
