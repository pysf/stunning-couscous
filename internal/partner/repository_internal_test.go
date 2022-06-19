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
			want: &Location{},
		},
		{
			args: args{
				point: `()`,
			},
			want: &Location{},
		},
		{
			args: args{
				point: ``,
			},
			want: &Location{},
		},
	}

	for i, tt := range tests {

		got := &Location{}
		err := got.ParsePostgresPoint(tt.args.point)
		if (err != nil) != tt.wantErr {
			t.Errorf("case(%d): ParsePostgresPoint() error = %v, wantErr %v", i, err, tt.wantErr)
			return
		}

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("case(%d): ParsePostgresPoint() = %v, want %v", i, got, tt.want)
		}

	}
}
