package tests

import (
	"net/http"
	"net/http/httptest"
	"news/aggregator"
	"news/db/models"
	"reflect"
	"testing"
)

func TestNewAggregator(t *testing.T) {
	tests := []struct {
		name string
		want aggregator.IAggregator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aggregator.NewAggregator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAggregator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAggregator_AggregateNews(t *testing.T) {
	type fields struct {
		client http.Client
	}
	type args struct {
		resource models.NewsResource
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.News
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := aggregator.Aggregator{
				Client: tt.fields.client,
			}
			got, err := a.AggregateNews(tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.AggregateNews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Aggregator.AggregateNews() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAggregator_aggregateFromHtml(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, _ = rw.Write([]byte(`
			<ul>
				<li class="srfrRow srfrRowIsOdd">
									<h3><a target="_blank" href="Test link">Test Title</a></h3>
									<span class="srfrFeedItemDate">Feb 26, 2021 | 17:01</span>
									<p>Test Description</p>
						
									<span class="srfrReadMore">
							<a target="_blank" href="https://news.drom.ru/Rolls-Royce-Cullinan-82962.html"></a>
						</span>
						
						<div class="clr"></div>
					</li>
			</ul>
			`))
	}))
	defer server.Close()
	client := http.Client{}
	serverRes, _ := client.Get(server.URL)
	type fields struct {
		client http.Client
	}
	type args struct {
		res      *http.Response
		resource models.NewsResource
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.News
		wantErr bool
	}{
		{
			name:   "Success parse",
			fields: fields{client: http.Client{}},
			args: args{
				res: serverRes,
				resource: models.NewsResource{
					IsRss: false,
					Link:  server.URL,
					ItemAttrs: ".srfrRow",
					PublishDateAttrs:".srfrFeedItemDate",
					TitleAttrs: "h3",
					DescriptionAttrs: "p",
					LinkAttrs: "a",
				},
			},
			want: []models.News{
				{
					PublishDate: "Feb 26, 2021 | 17:01",
					Title:       "Test Title",
					Link:        "Test link",
					Description: "Test Description",
					Resource: models.NewsResource{
						IsRss: false,
						Link:  server.URL,
						ItemAttrs: ".srfrRow",
						PublishDateAttrs:".srfrFeedItemDate",
						TitleAttrs: "h3",
						DescriptionAttrs: "p",
						LinkAttrs: "a",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := aggregator.Aggregator{
				Client: tt.fields.client,
			}
			got, err := a.AggregateFromHtml(tt.args.res, tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.aggregateFromHtml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Aggregator.aggregateFromHtml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAggregator_aggregateFromRss(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, _ = rw.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
				<rss version="2.0">
					<channel>
						<item>
							<guid>https://lenta.ru/news/2021/02/27/monuent/</guid>
							<author>Test Author</author>
  							<title>Test Title</title>
  							<link>Test link</link>
  							<description>Test Description</description>
  							<pubDate>31.01.1000</pubDate>
  							<category>Test Category</category>
						</item>
					</channel>
				</rss>
			`))
	}))
	defer server.Close()
	client := http.Client{}
	serverRes, _ := client.Get(server.URL)
	type fields struct {
		client http.Client
	}
	type args struct {
		res      *http.Response
		resource models.NewsResource
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.News
		wantErr bool
	}{
		{
			name:   "Success parse",
			fields: fields{client: http.Client{}},
			args: args{
				res: serverRes,
				resource: models.NewsResource{
					IsRss: true,
					Link:  server.URL,
				},
			},
			want: []models.News{
				{
					PublishDate: "31.01.1000",
					Title:       "Test Title",
					Category:    "Test Category",
					Author:      "Test Author",
					Link:        "Test link",
					Description: "Test Description",
					Resource: models.NewsResource{
						IsRss: true,
						Link:  server.URL,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := aggregator.Aggregator{
				Client: tt.fields.client,
			}
			got, err := a.AggregateFromRss(tt.args.res, tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.aggregateFromRss() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Aggregator.aggregateFromRss() = %v, want %v", got, tt.want)
			}
		})
	}
}
