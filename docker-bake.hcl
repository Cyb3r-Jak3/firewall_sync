target "alpine" {
  dockerfile = "./Dockerfile"
  target = "alpine"
}

target "distroless" {
  dockerfile = "./Dockerfile"
  target = "distroless"
}