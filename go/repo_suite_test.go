package main

import (
	"context"
	"log"
	"practice/pokeserver/testhelpers"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PokemonRepoTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  *Repository
	ctx         context.Context
}

func (suite *PokemonRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer
	repository, err := NewRepository(suite.ctx, suite.pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	suite.repository = repository
}

func (suite *PokemonRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

// Tear down your suite after each test
func (suite *PokemonRepoTestSuite) TearDownTest() {
	_, err := suite.repository.resetPokeVotes()
	if err != nil {
		log.Fatalf("Unable to reset sql table pokevotes: %s", err)
	}
}
func (suite *PokemonRepoTestSuite) TestCreatePokemon() {
	t := suite.T()

	pokemonCreated, err := suite.repository.createPokemonVote(Pokemon{
		Name: "Chari",
		Sprites: struct {
			BackDefault  string `json:"back_default"`
			FrontDefault string `json:"front_default"`
		}{
			FrontDefault: "url",
		},
		ID: 121,
	})
	assert.NoError(t, err)
	assert.True(t, pokemonCreated)
}

func (suite *PokemonRepoTestSuite) TestGetAllPokemon() {
	t := suite.T()
	suite.repository.createPokemonVote(Pokemon{
		Name: "Chari",
		Sprites: struct {
			BackDefault  string `json:"back_default"`
			FrontDefault string `json:"front_default"`
		}{
			FrontDefault: "url",
		},
		ID: 121,
	})
	allPokemon, err := suite.repository.getAllPokemonDBEntry()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(allPokemon))
}

func (suite *PokemonRepoTestSuite) TestGetPokemonById() {
	t := suite.T()
	_, err := suite.repository.createPokemonVote(Pokemon{
		Name: "Chari",
		Sprites: struct {
			BackDefault  string `json:"back_default"`
			FrontDefault string `json:"front_default"`
		}{
			FrontDefault: "url",
		},
		ID: 121,
	})
	assert.NoError(t, err)
	pokemon, err := suite.repository.getPokemonDBEntryById(121)
	testPokemon := PokeDBEntry{
		Id:   121,
		Name: "Chari",
		Vote: 0,
		Url:  "url",
	}
	assert.NoError(t, err)
	assert.NotNil(t, pokemon)
	assert.Equal(t, pokemon, testPokemon)
}

func (suite *PokemonRepoTestSuite) TestGetNonExistantPokemon() {
	t := suite.T()
	pokemon, err := suite.repository.getPokemonDBEntryById(-1)
	assert.Error(t, err)
	defaultPokeDbEntry := PokeDBEntry{
		Id:   0,
		Name: "",
		Vote: 0,
		Url:  "",
	}
	assert.Equal(t, pokemon, defaultPokeDbEntry)
}

func TestPokemonRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PokemonRepoTestSuite))
}
