module github.com/squarescale/hssh

require (
	github.com/alecthomas/gometalinter v3.0.0+incompatible // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/coreos/go-etcd v2.0.0+incompatible // indirect
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/golang/lint v0.0.0-20200302205851-738671d3881b // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/gordonklaus/ineffassign v0.0.0-20200809085317-e36bfde3bb78 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/manifoldco/promptui v0.8.0
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/nicksnyder/go-i18n v1.10.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/viper v1.7.1
	github.com/squarescale/cloudresolver v1.1.0
	github.com/squarescale/sshcommand v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/tsenart/deadcode v0.0.0-20160724212837-210d2dc333e9 // indirect
	github.com/ugorji/go/codec v0.0.0-20181204163529-d75b2dcb6bc8 // indirect
	github.com/xordataexchange/crypt v0.0.3-0.20170626215501-b2862e3d0a77 // indirect
	golang.org/x/build v0.0.0-20190111050920-041ab4dc3f9d // indirect
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897
	golang.org/x/tools v0.0.0-20201021122455-2be66b663cb6 // indirect
	google.golang.org/genproto v0.0.0-20201021134325-0d71844de594 // indirect
	gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20191105091915-95d230a53780 // indirect
)

go 1.14

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20200302205851-738671d3881b

replace golang.org/x/lint => github.com/golang/lint v0.0.0-20200302205851-738671d3881b
