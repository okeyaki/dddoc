# dddoc

dddoc is a tool for Golang which generates domain model documentation based on DDD.


## Install

```sh
go install github.com/okeyaki/dddoc/cmd/dddoc@latest
```


## Configure

```sh
curl -O https://raw.githubusercontent.com/okeyaki/dddoc/main/.dddoc.yml

vi .dddoc.yml
```


## Visualize your domain model

```sh
dddoc visualize
```

![](/doc/domain.png)


## Limitations

- Each component must have unique name
