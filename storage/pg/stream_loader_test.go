package pg

import (
	"context"
	"fmt"
	"testing"

	// "time"

	// "github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Test_streamLoader_Start(t *testing.T) {
	t.Run("running stream loader successfully", func(t *testing.T) {
		ctx := context.Background()
		req := testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{"5432/tcp"},
			Cmd:          []string{"postgres", "-c", "fsync=off"},
			Env:          map[string]string{"POSTGRES_PASSWORD": "postgres"},
			WaitingFor:   wait.ForLog(`PostgreSQL init process complete; ready for start up.`).AsRegexp(),
		}
		postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		testcontainers.CleanupContainer(t, postgresC)
		require.NoError(t, err)

		port, err := postgresC.MappedPort(ctx, "5432/tcp")
		require.NoError(t, err)

		pgURL := fmt.Sprintf("postgres://postgres:postgres@localhost:%d/postgres?sslmode=disable", port.Int())

		m, err := migrate.New("file://db/migrations", pgURL)
		require.NoError(t, err)

		require.NoError(t, m.Up())

		sl, err := NewStreamLoader(pgURL)
		require.NoError(t, err)

		sl.SetTable("products", "id").SetQueryLimit(4).SetSampleLimit(100)

		dataChan, errChan := sl.Start()

		var dataResult []map[string]any

	Loop:
		for {
			select {
			case data, ok := <-dataChan:
				if !ok {
					break Loop
				}
				dataResult = append(dataResult, data.(map[string]any))
			case err := <-errChan:
				require.NoError(t, err)
				break Loop
			}
		}

		assert.Len(t, dataResult, 16)
		assert.Equal(t, dataResult[0]["name"], "Item1")
		assert.Equal(t, dataResult[15]["name"], "Item16")
	})
}
