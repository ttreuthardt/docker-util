[![Build Status](https://travis-ci.org/ttreuthardt/docker-util.svg?branch=master)](https://travis-ci.org/ttreuthardt/docker-util)

# docker-util

Utility for generating simple config files at container startup.


## Config file
The json config file holds the required environment variables. Only
variables listed there can be used in templates. The templates array
defines all templates which shall be generated. fileMode, owner and
group are optional.

```json
{
  "envvars": [
    "MY_TEST_VAR"
  ],
  "templates": [
    {
      "templatePath": "/tests/test.tpl",
      "destPath": "/tests/dest/mytemplate.txt",
      "fileMode": "0700",
      "owner": "root",
      "group": "root"
    }
  ]
}
```


## Templates

Templates are simple go text template files. Environment variables are
accessible via the custom .ENV map.

```go

Value of env var MY_TEST_VAR={{ .Env.MY_TEST_VAR }}

```


## Usage

```sh
$ docker-util -config /path/to/config.json
```
