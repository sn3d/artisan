task "@go" "generate" {
	cmd = [ "go", "generate" ]
	deps = [ "//frontend:build" ]
}

task "@go" "build" {
    cmd = [ "go", "build" ]
    deps = [ ":generate" ]
}