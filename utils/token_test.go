package utils

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestNextToke(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "1",
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextToke()
			if (err != nil) != tt.wantErr {
				t.Errorf("NextToke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tokenPrefix := time.Now().Format("20060102150405")
			fmt.Printf("got : %s\n", got)
			if !strings.HasPrefix(got, tokenPrefix) {
				t.Errorf("NextToke() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenTimeOut(t *testing.T) {
	type args struct {
		token    string
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				token:    "20210609161500sdfasfa",
				duration: 1 * time.Hour,
			},
			want: false,
		},
		{
			name: "2",
			args: args{
				token:    "20210608161500sdfasfa",
				duration: 1 * time.Hour,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TokenTimeOut(tt.args.token, tt.args.duration); got != tt.want {
				t.Errorf("TokenTimeOut() = %v, want %v", got, tt.want)
			}
		})
	}
}
