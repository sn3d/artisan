task "go-1.13" "generate" {
    script = "go generate"
	deps = [
	    "//frontend:build"
	]
}

task "go-1.13" "build" {
    script = "go build"
    deps = [
        ":generate"
    ]
}