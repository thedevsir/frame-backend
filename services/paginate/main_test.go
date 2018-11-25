package paginate

import (
	"reflect"
	"testing"
)

func TestGenerate(t *testing.T) {

	type args struct {
		data  interface{}
		count int
		page  int
		limit int
	}
	tests := []struct {
		name string
		args args
		want *Paginate
	}{
		{
			name: "test",
			args: args{
				data: []string{
					"first",
					"secend",
				},
				count: 2,
				page:  1,
				limit: 10,
			},
			want: &Paginate{
				Data: []string{
					"first",
					"secend",
				},
				Pages: Pages{
					Current: 1,
					Prev:    0,
					HasPrev: false,
					Next:    2,
					HasNext: false,
					Total:   1,
				},
				Items: Items{
					Limit: 10,
					Begin: 1,
					End:   2,
					Total: 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Generate(tt.args.data, tt.args.count, tt.args.page, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
