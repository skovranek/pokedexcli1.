module github.com/skovranek/pokedexcli

go 1.14

replace (
	internal/pokeapi v1.0.0 => ./internal/pokeapi
	internal/pokecache v1.0.0 => ./internal/pokecache
)

require (
	internal/pokeapi v1.0.0
	internal/pokecache v1.0.0
)
