package main

import (
	"fmt"
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

var log logrus.Logger

func fallback() {
	syscall.Exec(viper.GetString("ssh"), os.Args, os.Environ())
}

func handleJump(args []string, provider string) []string {
	jh := viper.GetString(fmt.Sprintf("providers.%s.jumphost", provider))
	if jh == "" {
		log.Errorf("no jumphost specified")
		return args
	}
	r := cr.Resolvers[provider]
	hosts, err := r.Resolve(jh, viper.AllSettings())
	if err != nil {
		log.Errorf("error while resolving host")
		panic(err)
	}
	if len(hosts) == 0 {
		log.Errorf("resolution didn't returned any hosts")
		return args
	}

	ju := viper.GetString(fmt.Sprintf("providers.%s.jumpuser", provider))
	dest := hosts[0].Public
	if ju != "" {
		dest = fmt.Sprintf("%s@%s", ju, dest)
	}
	ssh_args_str := sshcommand.PrependOpt(args, []string{"-J", dest})
	log.Debugf("Adding ssh arguments: %#v", ssh_args_str)
	return ssh_args_str
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
		log.Errorf("error in prompt: %s", err)
	}
	return hosts[0]
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

	log.Debugf("starting hssh ...")

	if err != nil {
		log.Errorf("could not find config file")
	}

	ssh, err := exec.LookPath("ssh")
	if err != nil {
		panic("could not find ssh neither in path nor in configuration")
	}
	log.Debugf("Using ssh command found at: %#v", ssh)
	viper.SetDefault("ssh", ssh)

	provider := viper.GetString("provider")
	if provider == "" {
		log.Warnf("fallback: no provider specified")
		fallback()
	}

	args := os.Args
	sc, err := sshcommand.New(args)
	if err != nil {
		log.Warnf("fallback: ssh command not parseable with args: %s", os.Args)
		fallback()
	}
	desthost := sc.Hostname()
	if desthost != viper.GetString(fmt.Sprintf("providers.%s.jumphost", provider)) {
		args = handleJump(os.Args, provider)
	}
	args[0] = viper.GetString("ssh")

	r := cr.Resolvers[provider]

	hosts, err := r.Resolve(desthost, viper.AllSettings())
	if len(hosts) == 0 {
		log.Warnf("fallback: could not find any host matching destination %s", desthost)
		log.Warnf("%v", os.Args)
		log.Warnf("%v", viper.AllSettings())
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
