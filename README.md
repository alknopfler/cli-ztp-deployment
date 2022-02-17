# cli-ztp-deployment

# How to use it
```
Ztp is a command line to deploy ztp openshift clusters

Usage:
ztpcli [flags]
ztpcli [command]

Available Commands:
completion  Generate the autocompletion script for the specified shell
deploy      Commands to deploy things
help        Help about any command
mirror      Commands to mirroring
verify      Commands to verify things

Flags:
-h, --help   help for ztpcli

Use "ztpcli [command] --help" for more information about a command.
```  

# Verify

```
root:qct-d14u03 : ~/amorgant/cli-ztp-deployment {master}
$ ZTP_CONFIGFILE=./config.yaml ./cli-ztp-deployment verify -h
>>>> ConfigFile env is not empty. Reading file from this env
Commands to verify things

Usage:
  ztpcli verify [command]

Available Commands:
  httpd       Verify if File Server is running on the hub cluster
  mirror-ocp  Verify if the OCP mirring is successful based on mode (hub or spoke)
  mirror-ocp  Verify if the OLM operators are successful based on mode (hub or spoke)
  preflights  Run Preflight checks to validate the future deployments
  registry    Verify if registry is running on the server based on mode (hub | spoke)

Flags:
  -h, --help   help for verify

Use "ztpcli verify [command] --help" for more information about a command.
```

```
ZTP_CONFIGFILE=./config.yaml ./cli-ztp-deployment verify preflights

```

## Done by now
- preflights (verify) Done
- HTTPD (deployment and verify) Done
- registry (deployment and verify) Done
- Mirror ocp 