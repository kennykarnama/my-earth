## my-earth

Simple application to provide locations' weather

## Local development

### Bazel Setup

We use bazel as the build system. 

- Install bazelisk 
- Run bazel commands

```
bazelisk build //...
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