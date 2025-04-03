package formatx

import "testing"

func TestFormatSizePrecise(t *testing.T) {
	type args struct {
		bytes int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{
			100,
		}, "100 B"},
		{"", args{
			1024,
		}, "1.00 KB"},
		{"", args{
			1124,
		}, "1.10 KB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatSizePrecise(tt.args.bytes); got != tt.want {
				t.Errorf("FormatSizePrecise() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatSizePreciseWithUnits(t *testing.T) {
	type args struct {
		bytes int64
		unit  int64
		units []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{
			bytes: 1234,
			unit:  1000,
			units: defaultUnits,
		}, "1.23 KB"},
		{"", args{
			bytes: 1234,
			unit:  1000,
			units: []string{"X", "KX", "MX", "GX", "TX"},
		}, "1.23 KX"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatSizePreciseWithUnits(tt.args.bytes, tt.args.unit, tt.args.units); got != tt.want {
				t.Errorf("FormatSizePreciseWithUnits() = %v, want %v", got, tt.want)
			}
		})
	}
}
