workflow "Release" {
  on = "push"
  resolves = ["goreleaser"]
}


action "filter-tags" {
  uses = "actions/bin/filter@master"
  args = "tag"
  secrets = ["GITHUB_TOKEN"]
}

action "goreleaser" {
  uses = "docker://goreleaser/goreleaser"
  needs = ["filter-tags"]
  args = "release"
  secrets = [
    "GORELEASER_GITHUB_TOKEN",
  ]

}
