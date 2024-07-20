module github.com/n1haldev/distributed_file_storage

go 1.22.5

// require github.com/n1haldev/distributed_file_storage/p2p v0.0.0-20240720075946-21dcc34d6aee
replace github.com/n1haldev/distributed_file_storage/p2p => ./p2p

require (
	github.com/n1haldev/distributed_file_storage/p2p v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
