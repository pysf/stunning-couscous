package partner

import (
	"reflect"
	"testing"
)

func TestParsePostgresPoint(t *testing.T) {
	type args struct {
		point string
	}
	tests := []struct {
		args    args
		want    *Location
		wantErr bool
	}{
		{
			args: args{
				point: `(52.528971849007036,13.430548464498173)`,
			},
			want: &Location{
				Latitude:  52.528971849007036,
				Longitude: 13.430548464498173,
			},
		},
		{
			args: args{
				point: `(52.528971849007036,)`,
			},
			want: nil,
		},
		{
			args: args{
				point: `()`,
			},
			want: nil,
		},
		{
			args: args{
				point: ``,
			},
			want: nil,
		},
	}

	for _, tt := range tests {

		got := &Location{}
		err := got.ParsePostgresPoint(tt.args.point)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParsePostgresPoint() error = %v, wantErr %v", err, tt.wantErr)
			return
		}

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ParsePostgresPoint() = %v, want %v", got, tt.want)
		}

	}
}
