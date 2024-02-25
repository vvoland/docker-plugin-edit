group "default" {
    targets = ["bin-image-all"]
}

target "bin-image" {
    type = "image"
    platform = "local"
    tags = ["vlnd/docker-edit"]
}

target "bin-image-all" {
    inherits = ["bin-image"]
    platforms = [
        "linux/amd64",
        "linux/arm64",
        "darwin/amd64",
        "darwin/arm64",
        "windows/arm64",
        "windows/amd64"
    ]
}
