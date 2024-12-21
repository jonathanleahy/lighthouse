module github.com/jonathanleahy/project/webserver

go 1.21

require (
	github.com/jonathanleahy/project/jobscheduler v0.0.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/kr/pretty v0.3.1 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/jonathanleahy/project/jobscheduler => ../jobscheduler
