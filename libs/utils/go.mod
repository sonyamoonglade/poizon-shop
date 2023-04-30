module utils

go 1.18

require (
	github.com/sonyamoonglade/go_func v0.0.0-20230418180836-d7b9b025b11a
	github.com/stretchr/testify v1.8.2
)

require golang.org/x/exp v0.0.0-20220321173239-a90fa8a75705 // indirect

require (
	functools v0.0.0
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	functools => ../functools
	github.com/sonyamoonglade/go_func => ../../../go_func
)
