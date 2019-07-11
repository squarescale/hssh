package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	cr "github.com/squarescale/cloudresolver"
	"github.com/squarescale/sshcommand"
	"golang.org/x/crypto/ssh/terminal"
)

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
		log.Debugf(
			fmt.Sprintf(
				"Couldn't resolve host named \"%s\" with provider \"%s\", error: %s",
				jh,
				provider,
				err.Error(),
			),
		)
		return args
	}
	if len(hosts) == 0 {
		log.Debugf("resolution didn't returned any hosts")
		return args
	}

	ju := viper.GetString(fmt.Sprintf("providers.%s.jumpuser", provider))
	dest := hosts[0].Public
	if ju != "" {
		dest = fmt.Sprintf("%s@%s", ju, dest)
	}
	return sshcommand.PrependOpt(args, []string{"-J", dest})
}

func selectHost(hosts []cr.Host) cr.Host {
	if !terminal.IsTerminal(syscall.Stdin) {
		return hosts[0]
	}
	if viper.GetBool("interactive") {
		tmpls := promptui.SelectTemplates{
			Active:   `→  {{ .Id | cyan | bold }}`,
			Inactive: `   {{ .Id | cyan }}`,
			Selected: `{{ "✔" | green | bold }} {{ "Host" | bold }}: {{ .Id | cyan }}`,
			Details: `
provider: {{ .Provider   }}
region: {{ .Region     }}
zone: {{ .Zone       }}
id: {{ .Id         }}
private ipv4: {{ .PrivateIpv4}}
private ipv6: {{ .PrivateIpv6}}
private name: {{ .PrivateName}}
public ipv4: {{ .PublicIpv4 }}
public ipv6: {{ .PublicIpv6 }}
public name: {{ .PublicName }}`,
		}

		hostPrompt := promptui.Select{
			Label:     "Host",
			Items:     hosts,
			Templates: &tmpls,
		}

		idx, _, err := hostPrompt.Run()
		if err == nil {
			return hosts[idx]
		}
		log.Debugf("error in prompt: %s", err)
	}
	return hosts[0]
}

var log = logrus.New()

func main() {
	viper.SetConfigName("hssh")
	viper.AddConfigPath("$HOME/.config")
	viper.SetEnvPrefix("HSSH")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if viper.GetBool("debug") {
		log.SetLevel(logrus.DebugLevel)
	}

	logfn := viper.GetString("logfile")
	if logfn != "" {
		logfile, err := os.OpenFile(logfn, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("couldn't open logfile: %s: %s", logfn, err))
		}
		log.SetOutput(logfile)
	} else {
		if !terminal.IsTerminal(syscall.Stdout) {
			log.SetOutput(ioutil.Discard)
		}
	}

	log.Debugf("starting hssh ...")

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

	args := os.Args
	sc, err := sshcommand.New(args)
	if err != nil {
		log.Debugf("fallback: ssh command not parseable with args: %s", os.Args)
		fallback()
	}
	desthost := sc.Hostname()
	if desthost != viper.GetString(fmt.Sprintf("providers.%s.jumphost", provider)) {
		args = handleJump(os.Args, provider)
	}
	args[0] = viper.GetString("ssh")

	r := cr.Resolvers[provider]

	hosts, err := r.Resolve(desthost, viper.AllSettings())
	if err != nil {
		log.Debugf(
			fmt.Sprintf(
				"Couldn't resolve host named \"%s\" with provider \"%s\", error: %s",
				desthost,
				provider,
				err.Error(),
			),
		)
	}
	if len(hosts) == 0 {
		log.Debugf("fallback: could not find any host matching destination %s", desthost)
		log.Debugf("%v", os.Args)
		log.Debugf("%v", viper.AllSettings())
		fallback()
	}

	host := selectHost(hosts)
	hostname := host.Public
	if hostname == "" {
		hostname = host.Private
	}

	args = sshcommand.PrependOpt(args, []string{"-o", fmt.Sprintf("Hostname %s", hostname)})
	log.Debugf("executed command: %s", args)
	syscall.Exec(viper.GetString("ssh"), args, os.Environ())
}
