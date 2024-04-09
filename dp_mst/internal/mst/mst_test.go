package mst

import (
	cmn "dp_mst/internal/common"
	"reflect"
	"testing"
)

func TestKruskal(t *testing.T) {
	type args struct {
		root map[int64]*cmn.Graph
		mst  cmn.Graph
	}
	tests := []struct {
		name string
		args args
		want cmn.Graph
	}{
		{name: "paper",
			args: args{
				root: map[int64]*cmn.Graph{
					1: {{X: 1, Y: 2, W: 2.}, {X: 1, Y: 5, W: 1.}},
					2: {{X: 2, Y: 3, W: 1.}, {X: 2, Y: 5, W: 1.}, {X: 2, Y: 4, W: 7.}},
					3: {{X: 3, Y: 4, W: 4.}, {X: 3, Y: 5, W: 1.}},
				},
				mst: cmn.Graph{},
			},
			want: cmn.Graph{
				{X: 1, Y: 5, W: 1.},
				{X: 2, Y: 3, W: 1.},
				{X: 2, Y: 5, W: 1.},
				{X: 3, Y: 4, W: 4.},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make(cmn.Graph, 0)
			Kruskal(tt.args.root, tt.args.mst, cmn.Graph{}, &got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kruskal() = %v, want %v", got, tt.want)
			}
		})
	}
}
