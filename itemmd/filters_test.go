package main

import "testing"

func Test_isResourceNotItem(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "New York Journal (resource)",
			args: args{url: "https://www.loc.gov/resource/sn83030180/1897-04-11/ed-1/?sp=2"},
			want: true,
		},
		{
			name: "New York Journal (item)",
			args: args{url: "https://www.loc.gov/item/sn83030180/1897-04-11/ed-1/"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isResourceNotItem(tt.args.url); got != tt.want {
				t.Errorf("isResourceNotItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
