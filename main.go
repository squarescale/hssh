package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	cr "github.com/squarescale/cloudresolver"
	"github.com/squarescale/hssh/pkg/libhssh"
	"github.com/squarescale/sshcommand"
)

var (
	GitBranch        string
	GitCommit        string
	BuildDate        string
)

func init() {
	readline.Stdout = new(libhssh.NoBellStdOut)
}

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
	if r == nil {
		log.Errorf("fallback: no resolver found for provider: %s", provider)
		return args
	}
	hosts, err := r.Resolve(jh, viper.AllSettings())
	if err != nil {
		log.Errorf(
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
		log.Errorf("resolution didn't returned any hosts")
		return args
	}

	ju := viper.GetString(fmt.Sprintf("providers.%s.jumpuser", provider))
	dest := hosts[0].Public
	if ju != "" {
		dest = fmt.Sprintf("%s@%s", ju, dest)
	}
	// Check if ssh version supports -J as a valid option
	cmd := exec.Command("ssh", "-J")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	out, err := cmd.CombinedOutput()
	matched, err := regexp.MatchString("option requires an argument", string(out))
	if err != nil {
		log.Fatal(err)
	}
	var ssh_args_str []string

	opts := []string{
		"-o", fmt.Sprintf("ProxyCommand ssh -W %%h:%%p %s", dest),
	}

	if matched {
		opts = []string{"-J", dest}
	}

	ssh_args_str = sshcommand.PrependOpt(args, opts)

	log.Debugf("Adding ssh arguments: %#v", ssh_args_str)

	return ssh_args_str
}

func filterHosts(hosts []cr.Host, filter string) cr.Host {
	host := hosts[0]

	if filter == "" {
		return host
	}

	f, err := libhssh.NewFilterFromString(filter)
	if err != nil {
		log.Errorf("Bad filter: %s", err)
		return host
	}

	shosts := []cr.Host{}
	for _, host := range hosts {
		if f.HostMatch(&host) {
			shosts = append(shosts, host)
		}
	}

	if len(shosts) == 0 {
		log.Errorf("No host matching filter: %s", filter)
		os.Exit(2)
	}

	if len(shosts) > 1 {
		log.Errorf("Multiple hosts matching filter: %s", filter)
		os.Exit(3)
	}

	host = shosts[0]
	log.Debugf("Connecting to host: %#v", host)

	return host
}

func selectHost(hosts []cr.Host, filter string) cr.Host {
	if !terminal.IsTerminal(syscall.Stdin) {
		return filterHosts(hosts, filter)
	}
	if viper.GetBool("interactive") {
		tmpls := promptui.SelectTemplates{
			Active:   `→  {{ .Id | cyan | bold }}`,
			Inactive: `   {{ .Id | cyan }}`,
			Selected: `{{ "✔" | green | bold }} {{ "Host" | bold }}: {{ .Id | cyan }}`,
			Details: `
instance name: {{ .InstanceName }}
provider: {{ .Provider }}
region: {{ .Region }}
zone: {{ .Zone }}
id: {{ .Id }}
private ipv4: {{ .PrivateIpv4 }}
{{ if .PrivateIpv6 }}private ipv6: {{ .PrivateIpv6 }}
{{ end -}}
{{ if .PrivateName }}private name: {{ .PrivateName }}
{{ end -}}
{{ if .PublicIpv4 }}public ipv4: {{ .PublicIpv4 }}
{{ end -}}
{{ if .PublicIpv6 }}public ipv6: {{ .PublicIpv6 }}
{{ end -}}
{{ if .PublicName }}public name: {{ .PublicName }}
{{ end -}}
{{ if index .Tags "Project" }}project: {{ index .Tags "Project" }}
{{ end -}}
{{ if index .Tags "UserName" }}user: {{ index .Tags "UserName" }}
{{ end -}}
{{ if index .Tags "UserEmail" }}email: {{ index .Tags "UserEmail" }}
{{ end -}}`,
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
		os.Exit(1)
	}
	return filterHosts(hosts, filter)
}

var log = logrus.New()

func main() {
	viper.SetConfigName("hssh")
	viper.AddConfigPath("$HOME/.config")
	viper.SetEnvPrefix("HSSH")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("could not find config file")
	}

	if viper.GetBool("debug") {
		log.SetLevel(logrus.TraceLevel)
	}

	logfn := viper.GetString("logfile")
	if logfn != "" {
		logfile, err := os.OpenFile(logfn, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("couldn't open logfile: %q: %q", logfn, err)
		}
		log.SetOutput(logfile)
	} else {
		if !terminal.IsTerminal(syscall.Stdout) {
			log.SetOutput(ioutil.Discard)
		}
	}

	log.Debugf("starting hssh ...")

	ssh, err := exec.LookPath("ssh")
	if err != nil {
		log.Fatal("could not find ssh neither in path nor in configuration")
	}

	log.Debugf("Using ssh command found at: %q", ssh)
	viper.SetDefault("ssh", ssh)

	provider := viper.GetString("provider")
	if provider == "" {
		log.Warnf("fallback: no provider specified")
		fallback()
	}

	args := os.Args
	sc, err := sshcommand.New(args)
	if err != nil {
		if len(args) == 2 && args[1] == "-V" {
			fmt.Fprintf(os.Stderr, "hssh version: %s %s %s\n", GitBranch, GitCommit, BuildDate)
		} else {
			log.Warnf("fallback: ssh command not parseable with args: %s", os.Args)
		}
		fallback()
	}

	desthost := sc.Hostname()
	if desthost != viper.GetString(fmt.Sprintf("providers.%s.jumphost", provider)) {
		args = handleJump(os.Args, provider)
	}

	args[0] = viper.GetString("ssh")

	r := cr.Resolvers[provider]
	if r == nil {
		log.Warnf("fallback: no resolver found for provider: %s", provider)
		fallback()
	}

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
		log.Warnf("fallback: could not find any host matching destination %s", desthost)
		log.Warnf("%v", os.Args)
		log.Warnf("%v", viper.AllSettings())
		fallback()
	}

	host := selectHost(hosts, viper.GetString("filter"))
	hostname := host.Public
	if hostname == "" {
		hostname = host.Private
	}

	args = sshcommand.PrependOpt(args, []string{"-o", fmt.Sprintf("Hostname %s", hostname)})
	log.Debugf("executed command: %s", args)
	syscall.Exec(viper.GetString("ssh"), args, os.Environ())
}
