module dto

go 1.18

require (
	domain v0.0.0
	go.mongodb.org/mongo-driver v1.11.4
)

require functools v0.0.0 // indirect

replace (
	domain v0.0.0 => ../domain
	functools v0.0.0 => ../functools
)
