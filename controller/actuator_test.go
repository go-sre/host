package controller

import (
	"golang.org/x/time/rate"
	"net/url"
	"testing"
)

func TestParseLimitAndBurst(t *testing.T) {
	type args struct {
		values url.Values
	}
	tests := []struct {
		name    string
		args    args
		want    rate.Limit
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseLimitAndBurst(tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLimitAndBurst() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseLimitAndBurst() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseLimitAndBurst() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
