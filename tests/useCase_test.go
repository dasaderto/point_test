package tests

import (
	"news/db/models"
	"news/useCases"
	"sync"
	"testing"
)

func TestFind(t *testing.T) {
	type args struct {
		slice []string
		val   string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{
			name: "test success find",
			args: args{
				slice: []string{"test", "tester","testing"},
				val:   "test",
			},
			want: 0,
			want1: true,
		},
		{
			name: "test fail find",
			args: args{
				slice: []string{"test", "tester","testing"},
				val:   "testo",
			},
			want: -1,
			want1: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := useCases.Find(tt.args.slice, tt.args.val)
			if got != tt.want {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Find() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestActualizeNewsUC(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCases.ActualizeNewsUC()
		})
	}
}

func TestFetchResourceNews(t *testing.T) {
	type args struct {
		resource models.NewsResource
		workChan chan []models.News
		wg       *sync.WaitGroup
	}
	tests := []struct {
		name string
		args args
	}{

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCases.FetchResourceNews(tt.args.resource, tt.args.workChan, tt.args.wg)
		})
	}
}
