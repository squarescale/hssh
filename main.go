package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	cr "github.com/squarescale/cloudresolver"
	"github.com/squarescale/sshcommand"
)

var log logrus.Logger

func fallback() {
	syscall.Exec(viper.GetString("ssh"), os.Args, os.Environ())
}

func handleJump(args []string, provider string) []string {
	jh := viper.GetString(fmt.Sprintf("providers.%s.jumphost", provider))
	if jh == "" {
		log.Debugf("no jumphost specified")
		return args
	}
	r := cr.Resolvers[provider]
	hosts, err := r.Resolve(jh, viper.AllSettings())
	if err != nil {
		log.Debugf("error while resolving host")
		panic(err)
	}
	if len(hosts) == 0 {
		log.Debugf("resolution didn't returned any hosts")
		return args
	}

	ju := viper.GetString(fmt.Sprintf("providers.%s.jumpuser", provider))
	dest := ""
	if ju != "" {
		dest = fmt.Sprintf("%s@%s", ju, hosts[0].Public)
	} else {
		dest = hosts[0].Public
	}

	return sshcommand.PrependOpt(args, []string{"-J", dest})
}

func main() {
	log = *logrus.New()
	viper.SetConfigName("hssh")
	viper.AddConfigPath("$HOME/.config")
	viper.SetEnvPrefix("HSSH")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if viper.GetBool("debug") {
		log.SetLevel(logrus.DebugLevel)
	}
	if err != nil {
		log.Debugf("could not find config file")
	}

	ssh, err := exec.LookPath("ssh")
	if err != nil {
		panic("could not find ssh neither in path nor in configuration")
	}
	viper.SetDefault("ssh", ssh)

	provider := viper.GetString("provider")
	if provider == "" {
		log.Debugf("fallback: no provider specified")
		fallback()
	}

	args := handleJump(os.Args, provider)
	args[0] = viper.GetString("ssh")

	r := cr.Resolvers[provider]

	sc, err := sshcommand.New(args)
	if err != nil {
		log.Debugf("fallback: ssh command not parseable")
		fallback()
	}
	desthost := sc.Hostname()
	hosts, err := r.Resolve(desthost, viper.AllSettings())
	if len(hosts) == 0 {
		log.Debugf(fmt.Sprintf("fallback: could not find any host matching destination %s", desthost))
		fallback()
	}

	hostname := hosts[0].Public
	if hostname == "" {
		hostname = hosts[0].Private
	}

	args = sshcommand.PrependOpt(args, []string{"-o", fmt.Sprintf("Hostname %s", hostname)})
	syscall.Exec(viper.GetString("ssh"), args, os.Environ())
}
