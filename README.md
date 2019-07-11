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
brew install awscli && brew tap squarescale/sqssh && brew install sqssh
```

Please note also that to connect to the your infrastructure where you have set up some common SSH keys,
you will need to launch the following commands:

```
export SSH_AUTH_SOCK=/tmp/ssh-agent.sock
eval $(ssh-agent -a $SSH_AUTH_SOCK)
ssh-add your-privte-key
```

### Rebuilding

hssh can easily be rebuilt using Golang installed version (for instance using gvm https://github.com/moovweb/gvm):

```
gvm install go1.12.6 -b -B && \
git clone https://github.com/squarescale/hssh && \
cd hssh && \
go build .
```

The new hssh binary is located in the locally cloned git repository.

### Configuration file & usage

Here is an example of a ~/.config/hssh.yaml:

```
#debug: true
provider: aws
providers:
  aws:
    jumphost: bastion
    jumpuser: core
  gce:
    zone: europe-west1-b
```

Please note also that most of the variables defined there can be superseeded on the command line by uppercasing them and prefixing by HSSH_. For instance:

```
HSSH_DEBUG=1 AWS_PROFILE=dev HSSH_INTERACTIVE=1 hssh -o "StrictHostKeyChecking no" -o "UserKnownHostsFile /dev/null" core@nomad
```

Please note also the quotes around some of the standard ssh command line options which are required for hssh to properly pass them down to the underlying ssh command.

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
