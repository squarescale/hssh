module github.com/squarescale/hssh

require (
	github.com/chzyer/readline v1.5.1
	github.com/manifoldco/promptui v0.9.0
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/viper v1.13.0
	github.com/squarescale/cloudresolver v1.2.0
	github.com/squarescale/sshcommand v1.0.0
	github.com/stretchr/testify v1.8.0
	golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90
	golang.org/x/term v0.0.0-20220722155259-a9ba230a4035 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

go 1.16

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20200302205851-738671d3881b

replace golang.org/x/lint => github.com/golang/lint v0.0.0-20200302205851-738671d3881b
