package saveconfig

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveConfig(t *testing.T) {
	type args struct {
		filename string
		cfg      *Config
	}

	type test struct {
		args          args
		beforeSubTest func()
		afterSubTest  func()
		wantErr       error
	}

	tests := map[string]test{
		"Given valid config and filename Then return error nil": {
			args: args{
				filename: "config.json",
				cfg: &Config{
					Host: "127.0.0.1",
					Port: "8080",
				},
			},
			beforeSubTest: func() {},
			afterSubTest:  func() {},
		},
		"Given valid config and invalid filename When executed Then return error": {
			args: args{},
			beforeSubTest: func() {
				writeFile = func(filename string, data []byte, perm fs.FileMode) error {
					return errors.New("failed to write data to file")
				}
			},
			afterSubTest: func() {
				writeFile = os.WriteFile
			},
			wantErr: errors.New("failed to write data to file"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.beforeSubTest != nil {
				tt.beforeSubTest()
			}

			if tt.afterSubTest != nil {
				defer tt.afterSubTest()
			}

			args := tt.args
			err := SaveConfig(args.filename, args.cfg)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}
