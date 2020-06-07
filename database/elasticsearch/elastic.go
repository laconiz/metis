package elasticsearch

import (
	"context"
	"github.com/olivere/elastic/v7"
	"net/http"
	"time"
)

func New(option Option) (*Client, error) {

	var opts []elastic.ClientOptionFunc

	opts = append(opts, elastic.SetHttpClient(dialer))
	opts = append(opts, elastic.SetDecoder(&Decoder{}))
	opts = append(opts, elastic.SetURL(option.Address))
	opts = append(opts, elastic.SetSniff(option.Sniff))

	if option.User != "" && option.Password != "" {
		authFunc := elastic.SetBasicAuth(option.User, option.Password)
		opts = append(opts, authFunc)
	}

	client, err := elastic.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

type Option struct {
	Address  string // 地址
	User     string // 用户
	Password string // 密码
	Sniff    bool   // 嗅探
}

type Client struct {
	client *elastic.Client
}

// 创建索引
func (client *Client) Index(doc Document) error {

	exist, err := client.client.
		IndexExists(doc.ElkName()).
		Do(context.Background())

	if exist || err != nil {
		return err
	}

	_, err = client.client.
		CreateIndex(doc.ElkName()).
		BodyString(doc.ElkBody()).
		Do(context.Background())

	return err
}

func (client *Client) Insert(doc Document) error {

	_, err := client.client.
		Index().
		Index(doc.ElkName()).
		BodyJson(doc).
		Do(context.Background())

	return err
}

func (client *Client) InsertBulk(docs []Document) error {

	bulk := client.client.Bulk()
	for _, doc := range docs {

		request := elastic.NewBulkIndexRequest().
			Index(doc.ElkName()).
			Doc(doc)

		bulk.Add(request)
	}
	_, err := bulk.Do(context.Background())

	return err
}

var dialer = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 32,
		IdleConnTimeout:     time.Second * 10,
		MaxConnsPerHost:     32,
	},
	Timeout: time.Second * 10,
}
