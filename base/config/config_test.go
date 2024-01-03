package config

import (
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockFuncs struct {
	mock.Mock
}

func (m *mockFuncs) osStat(name string) (os.FileInfo, error) {
	args := m.Called(name)
	ret := args.Get(0).(os.FileInfo)
	e := args.Error(1)
	return ret, e
}

func (m *mockFuncs) osEnv() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *mockFuncs) osDirFS(dir string) fs.FS {
	args := m.Called(dir)
	return args.Get(0).(fs.FS)
}

type mockFS struct {
	mock.Mock
}

func (m *mockFS) ReadFile(name string) ([]byte, error) {
	args := m.Called(name)
	ret := args.Get(0).([]byte)
	e := args.Error(1)
	return ret, e
}

func (m *mockFS) Open(name string) (fs.File, error) {
	args := m.Called(name)
	ret := args.Get(0).(fs.File)
	e := args.Error(1)
	return ret, e
}

type mockFile struct {
	mock.Mock
}

func (f *mockFile) IsDir() bool {
	return f.Called().Bool(0)
}

func (f *mockFile) ModTime() time.Time {
	f.Called()
	return time.Now()
}

func (f *mockFile) Mode() fs.FileMode {
	f.Called()
	return 0
}

func (f *mockFile) Name() string {
	return f.Called().String(0)
}

func (f *mockFile) Size() int64 {
	f.Called()
	return 0
}

func (f *mockFile) Sys() interface{} {
	f.Called()
	return nil
}

type config struct {
	B second `namespace:"b"`
}

type second struct {
	B1 string   `env:"b1"`
	C  third    `namespace:"c"`
	B2 []string `env:"b2"`
}

type third struct {
	C1 string `env:"c1" required:"true"`
	C3 string `env:"c3" default:"env-default-c3-value"`
	C5 bool   `env:"c5" required:"true"`
}

type failConfig struct {
	B second `env:"b"`
}

type pointerConfig struct {
	A *string `env:"a"`
}

type parserFailConfig struct {
	A bool `env:"a"`
}

func TestParse(t *testing.T) {
	mf := new(mockFuncs)
	mfs := new(mockFS)
	osStat = mf.osStat
	osEnv = mf.osEnv
	osDirFS = mf.osDirFS

	tests := []struct {
		desc     string
		cfg      interface{}
		mockFunc func()
		exp      *config
		expErr   error
	}{
		{
			desc: "normal config",
			cfg:  &config{},
			mockFunc: func() {
				mf.On("osEnv").Return([]string{
					"b_b1=b_b1_value", "b_b2=b_b2_value1,b_b2_value2,b_b2_value3", "b_c_c1=b_c_c1_value", "b_c_c3=b_c_c3_value", "a1=yoyoyoyo", "b_c_c5=true",
				}).Once()
			},
			exp: &config{
				B: second{
					B1: "b_b1_value",
					B2: []string{"b_b2_value1", "b_b2_value2", "b_b2_value3"},
					C: third{
						C1: "b_c_c1_value",
						C3: "b_c_c3_value",
						C5: true,
					},
				},
			},
		},
		{
			desc: "use default value",
			cfg:  &config{},
			mockFunc: func() {
				mf.On("osEnv").Return([]string{
					"b_b1=b_b1_value", "b_c_c1=b_c_c1_value", "b_c_c5=true",
				}).Once()
			},
			exp: &config{
				B: second{
					B1: "b_b1_value",
					B2: []string{},
					C: third{
						C1: "b_c_c1_value",
						C3: "env-default-c3-value",
						C5: true,
					},
				},
			},
		},
		{
			desc: "required field",
			cfg:  &config{},
			mockFunc: func() {
				mf.On("osEnv").Return([]string{
					"b_b1=b_b1_value", "b_c_c3=b_c_c3_value",
				}).Once()
			},
			expErr: ErrValueRequired,
		},
		{
			desc: "struct but doesn't have namespace tag",
			cfg:  &failConfig{},
			mockFunc: func() {
				mf.On("osEnv").Return([]string{
					"b_b1=b_b1_value", "b_c_c1=b_c_c1_value", "b_c_c3=b_c_c3_value",
				}).Once()
			},
			expErr: ErrNamespaceTagNotFound,
		},
		{
			desc: "struct contains pointer field",
			cfg:  &pointerConfig{},
			mockFunc: func() {
				mf.On("osEnv").Return([]string{
					"b_b1=b_b1_value", "b_b2=b_b2_value1,b_b2_value2,b_b2_value3", "b_c_c1=b_c_c1_value", "b_c_c3=b_c_c3_value", "a1=yoyoyoyo",
				}).Once()
			},
			expErr: ErrInvalidType,
		},
		{
			desc: "parser parse failed",
			cfg:  &parserFailConfig{},
			mockFunc: func() {
				mf.On("osEnv").Return([]string{
					"a=asdf",
				}).Once()
			},
			expErr: ErrInvalidType,
		},
	}

	for _, test := range tests {
		envData = map[string]string{}
		test.mockFunc()

		err := Parse(test.cfg)
		if test.expErr != nil {
			assert.EqualError(t, test.expErr, err.Error(), test.desc)
		} else {
			assert.NoError(t, err, test.desc)
			assert.Equal(t, test.exp, test.cfg)
		}

		mf.AssertExpectations(t)
		mfs.AssertExpectations(t)
	}
}
