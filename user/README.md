### 课程中的 
git.imooc.com/`cap1573`/user 
需要更改为 
git.imooc.com/`coding-447`/user
# User Service

This is the User service

Generated with

```
micro new user
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- Alias: user

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./user
```

Build a docker image
```
make docker
```