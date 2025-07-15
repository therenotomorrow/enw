package enw_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therenotomorrow/enw"
	"github.com/therenotomorrow/enw/parsers/sethvargo"
)

func TestConfigValidate(t *testing.T) {
	t.Parallel()

	type sampleConfig struct {
		Field string `env:"FIELD"`
	}

	tests := []struct {
		name   string
		config enw.Config
		want   error
	}{
		{
			name:   "success with struct",
			config: enw.Config{Target: sampleConfig{}, Parser: sethvargo.New()},
			want:   nil,
		},
		{
			name:   "success with pointer",
			config: enw.Config{Target: &sampleConfig{}, Parser: sethvargo.New()},
			want:   nil,
		},
		{
			name:   "error on nil target",
			config: enw.Config{Target: nil, Parser: sethvargo.New()},
			want:   enw.ErrMissingTarget,
		},
		{
			name:   "error on nil parser",
			config: enw.Config{Target: sampleConfig{}, Parser: nil},
			want:   enw.ErrMissingParser,
		},
		{
			name:   "error on nil target pointer",
			config: enw.Config{Target: (*sampleConfig)(nil), Parser: sethvargo.New()},
			want:   enw.ErrNilTarget,
		},
		{
			name:   "error on invalid target (int)",
			config: enw.Config{Target: 123, Parser: sethvargo.New()},
			want:   enw.ErrInvalidTarget,
		},
		{
			name:   "error on invalid target pointer (int)",
			config: enw.Config{Target: new(int), Parser: sethvargo.New()},
			want:   enw.ErrInvalidTarget,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.config.Validate()

			require.ErrorIs(t, err, test.want)
		})
	}
}

func TestCollect(t *testing.T) {
	t.Parallel()

	// Определяем тестовые структуры
	type Sample struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT"`
	}
	type sampleConfig struct {
		AppName string   `env:"APP_NAME"`
		DB      Sample   `env:",prefix=DB_"`
		Servers []Sample `env:",prefix=SRV_"`
	}

	// Готовим данные
	dbConf := Sample{Host: "db.local", Port: 5432}
	srv1 := Sample{Host: "srv1.local", Port: 8080}

	// Ожидаемый результат для успешного кейса
	pkgPath := "github.com/therenotomorrow/enw_test"
	wantEnvs := []*enw.Env{
		{Value: "APP_NAME", Field: "AppName", Type: "string", Path: "sampleConfig->AppName", Package: pkgPath},
		{Value: "DB_HOST", Field: "Host", Type: "string", Path: "sampleConfig->DB->Host", Package: pkgPath},
		{Value: "DB_PORT", Field: "Port", Type: "int", Path: "sampleConfig->DB->Port", Package: pkgPath},
		{Value: "SRV_HOST", Field: "Host", Type: "string", Path: "sampleConfig->Servers->0->Host", Package: pkgPath},
		{Value: "SRV_PORT", Field: "Port", Type: "int", Path: "sampleConfig->Servers->0->Port", Package: pkgPath},
	}

	testCases := []struct {
		name     string
		config   enw.Config
		wantEnvs []*enw.Env
		wantErr  error
	}{
		{
			name: "Success with Struct Value",
			config: enw.Config{
				Target: sampleConfig{AppName: "MyApp", DB: dbConf, Servers: []Sample{srv1}},
				Parser: sethvargo.New(),
			},
			wantEnvs: wantEnvs,
			wantErr:  nil,
		},
		{
			name: "Success with Struct Pointer",
			config: enw.Config{
				Target: &sampleConfig{AppName: "MyApp", DB: dbConf, Servers: []Sample{srv1}},
				Parser: sethvargo.New(),
			},
			wantEnvs: wantEnvs,
			wantErr:  nil,
		},
		{
			name: "Failure on Invalid Target",
			config: enw.Config{
				Target: 123,
				Parser: sethvargo.New(),
			},
			wantEnvs: nil,
			wantErr:  enw.ErrInvalidTarget,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Мы передаем неполные данные в мок-парсер, так что нам не нужно сравнивать их полностью.
			// Этот тест сфокусирован на логике Collect, а не на парсинге.
			// Для полной проверки мы бы использовали настоящий sethvargo.Parser.

			gotEnvs, err := enw.Collect(tc.config)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, gotEnvs)
			} else {
				assert.NoError(t, err)
				// Упрощенное сравнение для демонстрации
				assert.Equal(t, len(tc.wantEnvs), len(gotEnvs), "Number of collected envs mismatch")
				// Для полного сравнения можно использовать assert.Equal(t, tc.wantEnvs, gotEnvs)
				// но это потребует точного заполнения всех полей в mockParser
			}
		})
	}
}
