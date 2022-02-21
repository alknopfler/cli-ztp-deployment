# cli-ztp-deployment

# How to use it
First of all, you will need the config.yaml (ztp_configfile) file with the information about your environment. You could take an example in this repo [/config.yaml](https://github.com/alknopfler/cli-ztp-deployment/blob/master/config.yaml)

After that, you need to export the next environment variables:
- KUBECONFIG: export KUBECONFIG=/path/to/kube-config.yaml
- ZTP_CONFIGFILE: export ZTP_CONFIGFILE=/path/to/config.yaml


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
$ ./ztpcli verify -h
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
ZTP_CONFIGFILE=./config.yaml ./ztpcli verify preflights
```

# Deploy

```
$ ./ztpcli deploy -h
>>>> ConfigFile env is not empty. Reading file from this env
Commands to deploy things

Usage:
  ztpcli deploy [command]

Available Commands:
  httpd       Deploy new File Server running on the hub cluster
  registry    Deploy new File Server running on the hub cluster based on mode (hub | spoke)

Flags:
  -h, --help   help for deploy

Use "ztpcli deploy [command] --help" for more information about a command.
```

```
$ ./ztpcli deploy registry -h
```

```
$ ./ztpcli deploy registry -h
>>>> ConfigFile env is not empty. Reading file from this env
Deploy new File Server running on the hub cluster based on mode (hub | spoke)

Usage:
  ztpcli deploy registry [flags]

Flags:
  -h, --help          help for registry
      --mode string   Mode of deployment for registry [hub|spoke]
```

```
./ztpcli deploy registry --mode hub
```



## Done by now
- preflights (verify) Done
- HTTPD (deployment and verify) Done
- registry (deployment and verify) Done
- Mirror ocp Done
- Mirror olm WIP (Podman push)