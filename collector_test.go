package enw_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/therenotomorrow/enw"
	"github.com/therenotomorrow/enw/parsers/sethvargo"
)

const (
	testPackage = "github.com/therenotomorrow/enw_test"
)

func TestNew(t *testing.T) {
	t.Parallel()

	obj := enw.New(sethvargo.New())

	assert.NotNil(t, obj)
}

func TestParsers(t *testing.T) {
	t.Parallel()

	var _ = []enw.Parser{
		&sethvargo.Parser{},
	}
}

func want() []*enw.Env {
	return []*enw.Env{
		{Value: "APP_NAME", Field: "AppName", Type: "string", Path: "AppName", Package: "", Tag: enw.Tag{Empty: true}},
		// Тип поля Host - string, а не Sample
		{Value: "DB_HOST", Field: "Host", Type: "string", Path: "DB->Host", Package: testPackage, Tag: enw.Tag{Empty: true}},
		{Value: "DB_PORT", Field: "Port", Type: "int", Path: "DB->Port", Package: testPackage, Tag: enw.Tag{Empty: true}},
		// То же самое для Cache
		{Value: "CACHE_HOST", Field: "Host", Type: "string", Path: "Cache->Host", Package: testPackage, Tag: enw.Tag{Empty: true}},
		{Value: "CACHE_PORT", Field: "Port", Type: "int", Path: "Cache->Port", Package: testPackage, Tag: enw.Tag{Empty: true}},
		// И для всех элементов срезов
		{Value: "SRV_HOST", Field: "Host", Type: "string", Path: "Servers->0->Host", Package: testPackage, Tag: enw.Tag{Empty: true}},
		{Value: "SRV_PORT", Field: "Port", Type: "int", Path: "Servers->0->Port", Package: testPackage, Tag: enw.Tag{Empty: true}},
		{Value: "SRV_HOST", Field: "Host", Type: "string", Path: "Servers->1->Host", Package: testPackage, Tag: enw.Tag{Empty: true}},
		{Value: "SRV_PORT", Field: "Port", Type: "int", Path: "Servers->1->Port", Package: testPackage, Tag: enw.Tag{Empty: true}},
		{Value: "PTR_SRV_HOST", Field: "Host", Type: "string", Path: "PtrServers->0->Host", Package: testPackage, Tag: enw.Tag{Empty: true}},
		{Value: "PTR_SRV_PORT", Field: "Port", Type: "int", Path: "PtrServers->0->Port", Package: testPackage, Tag: enw.Tag{Empty: true}},
		{Value: "PTR_SRV_HOST", Field: "Host", Type: "string", Path: "PtrServers->1->Host", Package: testPackage, Tag: enw.Tag{Empty: true}},
		{Value: "PTR_SRV_PORT", Field: "Port", Type: "int", Path: "PtrServers->1->Port", Package: testPackage, Tag: enw.Tag{Empty: true}},
		{Value: "NIL_SRV_HOST", Field: "Host", Type: "string", Path: "NilInSlice->0->Host", Package: testPackage, Tag: enw.Tag{Empty: true}},
		{Value: "NIL_SRV_PORT", Field: "Port", Type: "int", Path: "NilInSlice->0->Port", Package: testPackage, Tag: enw.Tag{Empty: true}},
		{Value: "NIL_SRV_HOST", Field: "Host", Type: "string", Path: "NilInSlice->2->Host", Package: testPackage, Tag: enw.Tag{Empty: true}},
		{Value: "NIL_SRV_PORT", Field: "Port", Type: "int", Path: "NilInSlice->2->Port", Package: testPackage, Tag: enw.Tag{Empty: true}},
	}

}

func TestCollectorCollect(t *testing.T) {
	t.Parallel()

	type Sample struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT"`
	}

	type sampleConfig struct {
		AppName    string    `env:"APP_NAME"`
		DB         Sample    `env:",prefix=DB_"`
		Cache      *Sample   `env:",prefix=CACHE_"`
		EmptyCache *Sample   `env:",prefix=EMPTY_"`
		Servers    []Sample  `env:",prefix=SRV_"`
		PtrServers []*Sample `env:",prefix=PTR_SRV_"`
		NilInSlice []*Sample `env:",prefix=NIL_SRV_"`
		unexported string    `env:"UNEXPORTED"`
	}

	var (
		srv1      = Sample{Host: "srv1.local", Port: 8080}
		srv2      = Sample{Host: "srv2.local", Port: 8081}
		dbConf    = Sample{Host: "db.local", Port: 5432}
		cacheConf = Sample{Host: "cache.local", Port: 6379}
	)

	tests := []struct {
		name  string
		input any
		want  []*enw.Env
	}{
		{
			name: "full config walk as a struct",
			input: sampleConfig{
				AppName:    "MyApp",
				DB:         dbConf,
				Cache:      &cacheConf,
				EmptyCache: nil,
				Servers:    []Sample{srv1, srv2},
				PtrServers: []*Sample{&srv1, &srv2},
				NilInSlice: []*Sample{&srv1, nil, &srv2},
			},
			want: want(),
		},
		{
			name: "full config walk as a pointer",
			input: &sampleConfig{
				AppName:    "MyApp",
				DB:         dbConf,
				Cache:      &cacheConf,
				EmptyCache: nil,
				Servers:    []Sample{srv1, srv2},
				PtrServers: []*Sample{&srv1, &srv2},
				NilInSlice: []*Sample{&srv1, nil, &srv2},
			},
			want: want(),
		},
		{name: "nil struct", input: nil, want: []*enw.Env{}},
		{name: "empty struct", input: struct{}{}, want: []*enw.Env{}},
		{name: "not a struct", input: 123, want: []*enw.Env{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			obj := enw.New(sethvargo.New())

			got := obj.Collect(reflect.ValueOf(test.input), "", "", "")

			assert.Len(t, got, len(test.want))
			assert.Equal(t, test.want, got)
		})
	}
}
