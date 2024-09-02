package main

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRepository(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithInitScripts(filepath.Join("testdata", "init-db.sql")),
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("mypostgresuser"),
		postgres.WithPassword("mysecretpassword"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	pokeRepo, err := NewRepository(ctx, connStr)

	assert.NoError(t, err)

	pokemon := Pokemon{
		Name: "Chari",
		Sprites: struct {
			BackDefault  string `json:"back_default"`
			FrontDefault string `json:"front_default"`
		}{
			FrontDefault: "url",
		},
		ID: 121,
	}
	p := pokeRepo.createPokemonVote(pokemon)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.True(t, p)
}
