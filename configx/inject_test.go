package configx

import "testing"

func TestGetWorkDirectoryKey(t *testing.T) {
	type args struct {
		workDir string
	}
	tests := []struct {
		name    string
		args    args
		wantKey string
		wantErr bool
	}{
		{
			"",
			args{
				workDir: "d:\\a\\bb",
			},
			"d\\a\\bb",
			false,
		},
		{
			"",
			args{
				workDir: "e:",
			},
			"e",
			false,
		},
		{
			"",
			args{
				workDir: "",
			},
			"",
			true,
		},
		{
			"",
			args{
				workDir: "/",
			},
			"",
			false,
		},
		{
			"",
			args{
				workDir: "/a",
			},
			"a",
			false,
		},
		{
			"",
			args{
				workDir: "/ab/cd",
			},
			"ab/cd",
			false,
		},
		{
			"",
			args{
				workDir: "/ab\\cd",
			},
			"ab\\cd",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, err := GetWorkDirectoryKey(tt.args.workDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWorkDirectoryKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotKey != tt.wantKey {
				t.Errorf("GetWorkDirectoryKey() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func TestGetDefaultAppConfigRoot(_ *testing.T) {
	_, _ = GetDefaultAppConfigRoot()
}

func TestGetWorkDirectoryKeyEx(t *testing.T) {
	type args struct {
		workDir    string
		trimPrefix []string
	}
	tests := []struct {
		name    string
		args    args
		wantKey string
		wantErr bool
	}{
		{
			"",
			args{
				workDir:    "E:\\work",
				trimPrefix: []string{"E:\\"},
			},
			"work",
			false,
		},
		{
			"",
			args{
				workDir:    "/home/z/work/a11",
				trimPrefix: []string{"/home/z/"},
			},
			"work/a11",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, err := GetWorkDirectoryKeyEx(tt.args.workDir, tt.args.trimPrefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWorkDirectoryKeyEx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotKey != tt.wantKey {
				t.Errorf("GetWorkDirectoryKeyEx() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}
