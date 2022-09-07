Hype/Heaven SSH wrapper for Cloud infrastructures
=================================================

Usage
-----

### Overview

When used with a AWS infrastructure, where EC2 instances are named, hssh  allows you to connect via a jumphost
directly to non publically accessible instances.

### Installation

You can download official binaries from https://github.com/squarescale/hssh/releases

On MacOSX, use the following:

```
brew install awscli && brew tap squarescale/hssh && brew install hssh
```

### Common SSH options and basic usage

as hssh can be seen as an SSH wrapper, all standard SSH options can be used on the hssh command line. As an example:

```
hssh -o "StrictHostKeyChecking no" -o "UserKnownHostsFile /dev/null" core@*
```

Note that the star (*) as the host will be checked for matches in cloudresolver library which will connect to your Cloud infrastructure (depending on the value of HSSH_PROVIDER environment variable) to retrieve potential matches using instances names, tags, IP adresses and more.

### Important note about SSH preset keys and ssh-agent

Please note also that to connect to your infrastructure where you have set up some common SSH keys though a `bastion`,
you will need to launch the following commands:

```
export SSH_AUTH_SOCK=/tmp/ssh-agent.sock
eval $(ssh-agent -a $SSH_AUTH_SOCK)
ssh-add your-private-key
```

Specifying `-i you-private-key` file as a command line option will only work for the `bastion` host but afterwards, using `bastion` as a jumphost will fail

### Rebuilding

hssh can easily be rebuilt using Golang installed version (for instance using gvm https://github.com/moovweb/gvm):

```
gvm install $(gvm listall | grep go1.19 | tail -1) -b -B && \
gvm use $(gvm listall | grep go1.19 | tail -1) && \
git clone https://github.com/squarescale/hssh && \
cd hssh && make
```

The new hssh binary is located in the locally cloned git repository.

### Debugging locally

When you want to rebuild and debug **hssh** code locally, you might want to also use a local version of the most usefull dependancies like **cloudresolver** and **sshcommand**

In order to do so, just add the following line to the end of go.mod file:

```
replace github.com/squarescale/cloudresolver => ../cloudresolver
```

### Making changes in modules used by hssh

You can get versions of Go modules using the following command:

```go list [-json] -m -versions github.com/squarescale/cloudresolver```

If you create a test branch for a module, you can retrieve the corresponding tag to insert in go.mod
by using the following command:

```go list -m github.com/squarescale/cloudresolver@<branch-name>```

You then need to update go.mod file accordingly before rebuilding.

To list all dependencies modules versions:

```
go list -m -u all
```

To get the latest pseudo-version of a module, go into your local version of the module git repository and type:

```
TZ=UTC git --no-pager show \
  --quiet \
  --abbrev=12 \
  --date='format-local:%Y%m%d%H%M%S' \
  --format="%cd-%h"
```

This should give an output like:

```
20200630191459-0b5bf24f1853
```

which you can then use in go.mod

### Configuration file & usage

Here is an example of a ~/.config/hssh.yaml:

```
#debug: true
provider: aws
providers:
  aws:
    jumphost: bastion
    jumpuser: core
  azure:
    jumphost: bastion
    jumpuser: debug
  gce:
    zone: europe-west1-b
```

Please note also that most of the variables defined there can be superseeded on the command-line by uppercasing them and prefixing by HSSH_. For instance:

```
HSSH_DEBUG=1 AWS_PROFILE=dev HSSH_INTERACTIVE=1 hssh -o "StrictHostKeyChecking no" -o "UserKnownHostsFile /dev/null" core@nomad
```

Please also note the quotes around some of the standard ssh command-line options which are required for hssh to properly pass them down to the underlying ssh command.

### Advanced usage

You can easily use any command using SSH with hssh. For instance to rsync a file/dir from any infrastructure node to your local host, you can use:

```
AWS_PROFILE=dev rsync -e 'hssh -o "StrictHostKeyChecking no" -o "UserKnownHostsFile /dev/null"' -av core@nomad:/etc .
```
which is also equivalent to:
```
AWS_PROFILE=dev RSYNC_RSH='hssh -o "StrictHostKeyChecking no" -o "UserKnownHostsFile /dev/null"' rsync  -av core@nomad:/etc .
```

A logfile option is also provided if you want to redirect log messages into a specific file. If this option is not used and
stdout of hssh is redirected (aka not isatty), then log messages are discarded.

### Azure environment setup

To be able to target Azure platform (added in tag v0.1.6 of hssh and v1.0.9 of cloudresolver), some specific operations have to be carried out.

However, if you do not want
to follow the steps describe hereafter and if you have a working az CLI configuration, you can use it by swapping
the comment tags of lines 31 and 32 in [Cloudresolver azure.go](https://github.com/squarescale/cloudresolver/blob/master/azure.go#L31)

In order to use Azure Golang SDK, you need to configure your Azure account properly so that the SDK calls can succeed. This is, at the time of this writing, not very well documented and a lot of people seem to have had problems putting this into working mode.

Thanks to some people in the [Azure SDK for Go Gophers slack channel](https://gophers.slack.com/archives/CA7HK8EEP), we have been able to quickly set up a working environment.

The following set of environment variables should be specified when using Azure provider for cloudresolver library:

```
AZURE_CLIENT_ID
AZURE_CLIENT_SECRET
AZURE_LOCATION_DEFAULT
AZURE_SUBSCRIPTION_ID
AZURE_TENANT_ID
```

`AZURE_TENANT_ID` and `AZURE_SUBSCRIPTION_ID` are the same than in the output you get when running `az login` command.

In our case, `AZURE_LOCATION_DEFAULT` is set to `westeurope`.

In order to properly set the value of `AZURE_CLIENT_ID ` and `AZURE_CLIENT_SECRET`, you need to configure a RBAC 'role' running the following az CLI command:

```
az ad sp create-for-rbac
```

The output should look like:

```
{
  "appId": "<SOME_APP_UUID_TO_BE_USED_AS_AZURE_CLIENT_ID_ENV_VAR>",
  "displayName": "azure-cli-2020-06-17-11-59-33",
  "name": "http://azure-cli-2020-06-17-11-59-33",
  "password": "<THE_PASSWORD_TO_BE_USED_AS_AZURE_CLIENT_SECRET_ENV_VAR>",
  "tenant": "<YOUR_TENANT_ID_SAME_AS_AZURE_TENANT_ID_ENV_VAR>"
}
```

Afterwards, you can use `az ad sp list` to retrieve existing RBAC entries in your Azure account and also you can check validity of the password and ID using:

```
az login --service-principal --username $AZURE_CLIENT_ID --password $AZURE_CLIENT_SECRET --tenant $AZURE_TENANT_ID
```
