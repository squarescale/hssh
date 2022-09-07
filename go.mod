module github.com/squarescale/hssh

require (
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/manifoldco/promptui v0.8.0
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/viper v1.7.1
	github.com/squarescale/cloudresolver v1.1.0
	github.com/squarescale/sshcommand v1.0.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897
	google.golang.org/genproto v0.0.0-20201021134325-0d71844de594 // indirect
)

go 1.16

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20200302205851-738671d3881b

replace golang.org/x/lint => github.com/golang/lint v0.0.0-20200302205851-738671d3881b
