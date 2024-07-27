## my-earth

Simple application to provide locations' weather

## Local development

### Bazel Setup

We use bazel as the build system. 

- Install bazelisk
- Install bazel dependencies
  - python3
  - cclang toolchains, just to be safe, please install build-essentials
- Run bazel commands

```
bazelisk build //...
```

Run earth-server

```
bazelisk run //cmd/earth-server:earth-server
```

Running bazelisk build will create binary of each of go package.

It also creates local docker image for:

+ earth-server
+ earth-migrator

To push the image

```
bazelisk run //cmd/earth-server:push
```

```
bazelisk run //src/infra/db/migrations:push 
```