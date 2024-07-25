# BUILD
# gazelle:./api/openapi/genapi
load("@bazel_gazelle//:def.bzl", "gazelle")

gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune",
    ],
    command = "update-repos",
)

gazelle(name = "gazelle")
