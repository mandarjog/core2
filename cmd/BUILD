load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["cmd.go"],
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_golang_glog//:go_default_library",
        "@com_github_hashicorp_go_multierror//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@io_istio_api//:proxy/v1/config",
        "@io_istio_core//proxy:go_default_library",
        "@io_istio_core//tools/version:go_default_library",
    ],
)
