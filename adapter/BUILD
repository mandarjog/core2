# gazelle:ignore
package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load(":inventory.bzl", "inventory_library")

inventory_library(
    name = "go_default_library",
    packages = {
        "denier": "istio.io/core/adapter/denier",
        "list": "istio.io/core/adapter/list",
        "noop": "istio.io/core/adapter/noop",
        "prometheus": "istio.io/core/adapter/prometheus",
        "stackdriver": "istio.io/core/adapter/stackdriver",
        "statsd": "istio.io/core/adapter/statsd",
        "stdio": "istio.io/core/adapter/stdio",
        "svcctrl": "istio.io/core/adapter/svcctrl",
        "memquota": "istio.io/core/adapter/memquota",
    },
    deps = [
        "//adapter/denier:go_default_library",
        "//adapter/list:go_default_library",
        "//adapter/memquota:go_default_library",
        "//adapter/noop:go_default_library",
        "//adapter/prometheus:go_default_library",
        "//adapter/stackdriver:go_default_library",
        "//adapter/statsd:go_default_library",
        "//adapter/stdio:go_default_library",
        "//adapter/svcctrl:go_default_library",
    ],
)
