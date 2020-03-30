# Onix Terraform Provider <img src="../docs/pics/ox.png" width="200" height="200" align="right">

The Terraform provider for Onix allows Terraform to interact with Onix Web API.

The provider documentation can be found [here](docs/index.md).

[GNU Make](https://www.gnu.org/software/make/) is used to build, test and install the provider hiding the details of how this is done from the user.

The following table shows the [Makefile](Makefile) commands available to the user:

| Command | Description |
|---|---|
| `test` | *Runs the provider acceptance tests against a containerised backend Onix Web API / database. It requires a host where docker is installed. <br> To make things easy, tests run using the default credentials: `admin:0n1x`. <br> When running the tests against backends other than the containerised test services, the credentials used by the provider to connect to the Web API can be overriden using the environment variables: TF_PROVIDER_OX_USER, TF_PROVIDER_OX_PWD, and TF_PROVIDER_OX_URI.* |
| `install` | *Compiles, installs and initialises the Terraform provider for Onix.* |
| `build` | *Builds the terraform provider for the current OS.* |
| `package` | *Builds the terraform providers for Linux, Windows and Darwin OS respectively. Each provider is zipped and placed in the **build** folder.* |
| `package_linux` | *Builds the terraform providers for Linux.* |
| `package_windows` | *Builds the terraform providers for Windows.* |
| `package_darwin` | *Builds the terraform providers for MacOS.* |

For example, to run the acceptance tests simply run:

```bash
$ make test
```

## Required tools

- [GNU Make](https://www.gnu.org/software/make/)
- [Golang](https://golang.org/) 
- [Terraform](https://www.terraform.io/)

## Installing the provider

For information on how to install this provider see [here](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).


