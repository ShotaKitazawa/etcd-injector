etcd-injector
===

![License](https://img.shields.io/github/license/ShotaKitazawa/etcd-injector)
![test](https://github.com/ShotaKitazawa/etcd-injector/workflows/test/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/ShotaKitazawa/etcd-injector)
![Dependabot](https://badgen.net/dependabot/ShotaKitazawa/etcd-injector?icon=dependabot)

recursive copy & inject (or replace) json value of etcd

![comcept](./images/comcept.png)

## Install

* download binary from [GitHub Release](https://github.com/ShotaKitazawa/etcd-injector/releases)

## Usage

```
USAGE:
   etcd-injector [options]

OPTIONS:
   --src-endpoints value             source endpoints of etcd [$ETCD_SRC_ENDPOINTS]
   --src-username value              username of source etcd [$ETCD_SRC_USERNAME]
   --src-password value              password of source etcd [$ETCD_SRC_PASSWORD]
   --src-directory value, -s value   source directory of etcd [$ETCD_SRC_DIRECTORY]
   --dst-endpoints value             destination endpoints of etcd [$ETCD_DST_ENDPOINTS]
   --dst-username value              username of destination etcd [$ETCD_DST_USERNAME]
   --dst-password value              password of destination etcd [$ETCD_DST_PASSWORD]
   --dst-directory value, -d value   destination directory of etcd [$ETCD_DST_DIRECTORY]
   --rules-filepath value, -f value  path of file written injection rules [$RULES_FILEPATH]
   --ignore value                    specified "--ignore=/key", "xxx" is excluded from copy target [$IGNORE_KEYS]
   --delete                          delete dst key if does not exist in src (like "rsync --delete") (default: false)
   --verbose, -x                     output results of replacement (default: false)
   --help, -h                        show help (default: false)
   --version, -v                     print the version (default: false)
```

### Example

```
etcd-injector \
  --src-endpoints=http://127.0.0.1:2379 \
  --dst-endpoints=http://127.0.0.1:2379 \
  --src-directory=/src/com/example \
  --ignore=/skydns/com/example/local \
  --dst-directory=/dst/com/example \
  --rules-filepath=./example/rules.yaml \
  --delete
```

* copied `/src/com/example` to `/dst/com/example` in `http://127.0.0.1:2379`
* injected `{"value": "replaced"}` in value of `/dst/com/example/*` by ./example/rules.yaml
* ignored to copy `/src/com/example/local`

