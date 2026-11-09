module github.com/heliantheon/aegis-go/service

go 1.26

require (
	aidanwoods.dev/go-paseto v1.6.0
	github.com/go-json-experiment/json v0.0.0-20260214004413-d219187c3433
	github.com/heliantheon/aegis-go/utilities v0.0.2
)

require (
	aidanwoods.dev/go-result v0.3.1 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
)

replace github.com/heliantheon/aegis-go/utilities => ../utilities
