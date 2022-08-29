module audit

go 1.19

// using jbsmith7741's fork with duration humanization features
replace github.com/dustin/go-humanize => github.com/jbsmith7741/go-humanize v1.0.1-0.20211011174707-9d50e1685b88

require (
	github.com/deckarep/golang-set/v2 v2.1.0
	github.com/diamondburned/arikawa/v3 v3.0.1-0.20220822214349-9e9f90a65248
	github.com/rs/zerolog v1.27.0
	golang.org/x/exp v0.0.0-20220826205824-bd9bcdd0b820
)

require (
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
)
