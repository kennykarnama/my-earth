# gazelle:ignore

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//:build/oapi-codegen.bzl", "oapi_codegen")

oapi_codegen(
    name = "api",
    out = "meteo_source.gen.go",
    generate = [
        "types",
        "client",
    ],
    package = "genapi",
    spec = "meteor_source.yaml",
    visibility = ["//visibility:public"],
)

oapi_codegen(
    name = "location_api",
    out = "location.gen.go",
    generate = [
        "types",
        "gin-server",
    ],
    package = "genapi",
    spec = "location.yaml",
    visibility = ["//visibility:public"],
)

go_library(
    name = "openapi",
    srcs = [
        "location.gen.go",
        "meteo_source.gen.go",
    ],
    importpath = "github.com/kennykarnama/my-earth/api/openapi",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_gin_gonic_gin//:gin",
        "@com_github_oapi_codegen_runtime//:runtime",
        "@com_github_oapi_codegen_runtime//types",
    ],
)
