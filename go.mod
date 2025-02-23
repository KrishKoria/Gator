module github.com/KrishKoria/Gator

go 1.23.5

replace github.com/KrishKoria/GatorConfig v0.0.0 => ./internal/config/

require github.com/KrishKoria/GatorConfig v0.0.0

require (
	github.com/google/uuid v1.6.0 
	github.com/lib/pq v1.10.9 
)
