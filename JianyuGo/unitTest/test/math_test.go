/**
* @Author:zhoutao
* @Date:2022/1/29 10:10
* @Desc:
 */

package math

import "testing"

func TestAdd(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"negative+negative",
			args{-1, -1},
			-2,
		},
		{"negative+positive",
			args{-1, 1},
			0,
		},
		{"positive+positive",
			args{1, 1},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}
