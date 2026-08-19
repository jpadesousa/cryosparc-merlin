<<<<<<< HEAD
# CryoSPARC Merlin

A command-line helper for installing, configuring, and managing [CryoSPARC](https://cryosparc.com/) on the Merlin cluster.

The tool provides a single CLI for common CryoSPARC administration tasks, including master/worker installation, lane management, user creation, and direct access to `cryosparcm` and `cryosparcw` commands.

## Features

* Install a complete CryoSPARC deployment
* Install multiple versions of CryoSPARC
* Install CryoSPARC master and worker nodes independently
* Support `x86_64` and `aarch64` architectures
* Create and manage CryoSPARC lanes
* Install and remove lanes
* Create CryoSPARC users
* Run `cryosparcm` commands
* Run `cryosparcw` commands
* Automatic detection of available ports with [findbaseport](https://github.com/jpadesousa/findbaseport)
* Detection of running CryoSPARC instances with [uports](https://github.com/jpadesousa/uports)

## Installation

Clone the repository and build the executable:

```bash
git clone <repository-url>
cd cryosparc-merlin

make build
```

Verify the installation:

```bash
cryosparc-merlin --help
```

To display the installed version:

```bash
cryosparc-merlin --version
```

## Usage

The general command structure is:

```text
cryosparc-merlin <command> [flags] [arguments]
```

Available top-level commands:

```text
cryosparcm
cryosparcw
install
instances
lanes
user
```

Run the following to see the available commands and flags:

```bash
cryosparc-merlin --help
```

By default, CryoSPARC will be installed in `$HOME/cryosparc`. However, this can
be changed with the `-d/--cryosparc-path`:
```bash
cryosparc-merlin install complete --license [LICENSE] --version 5.0.6 -d $HOME/cryosparc_v5.0.6
```

When using the `cryosparc-merlin` commands that will target this installation, add `-d $HOME/cryosparc_v5.0.6` to the command. For example, to create and install all default lanes in this installation:

```bash
cryosparc-merlin lanes create default -d $HOME/cryosparc_v5.0.6
```

and then install all lanes created
```bash
cryosparc-merlin lanes install -d $HOME/cryosparc_v5.0.6 --all
```
=======
# cryosparc-merlin

A CLI helper for installing, configuring, and managing CryoSPARC on the PSI Merlin cluster.
>>>>>>> f0cfa93 (Initial commit)
