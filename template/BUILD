package(default_visibility = ["//visibility:public"])

load("//tools/codegen:generate.bzl", "mixer_supported_template_library")

mixer_supported_template_library(
    name = "go_default_library",
    packages = {
        "//template/listentry:go_default_library_proto.descriptor_set": "istio.io/core/template/listentry",
        "//template/logentry:go_default_library_proto.descriptor_set": "istio.io/core/template/logentry",
        "//template/metric:go_default_library_proto.descriptor_set": "istio.io/core/template/metric",
        "//template/quota:go_default_library_proto.descriptor_set": "istio.io/core/template/quota",
        "//template/reportnothing:go_default_library_proto.descriptor_set": "istio.io/core/template/reportnothing",
        "//template/checknothing:go_default_library_proto.descriptor_set": "istio.io/core/template/checknothing",
    },
    deps = [
        "//template/checknothing:go_default_library",
        "//template/listentry:go_default_library",
        "//template/logentry:go_default_library",
        "//template/metric:go_default_library",
        "//template/quota:go_default_library",
        "//template/reportnothing:go_default_library",
    ],
)
