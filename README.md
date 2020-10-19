# Artisan

[![Release](https://img.shields.io/github/release/unravela/artisan.svg?style=flat-square)](https://github.com/unravela/artisan/releases/latest)
[![Build](https://img.shields.io/github/workflow/status/unravela/artisan/build?style=flat-square)](https://github.com/unravela/artisan/actions?query=workflow%3Abuild)


Artisan is a build orchestrator for mono repositories powered by Docker. Artisan 
helps you build complex codebases without the need to install all build tools.
 
## How it works
Let's have a repository with an application written in Java and build by Gradle 
and Vue frontend. Usually, we need to install the correct version of NPM, Java, 
and Gradle. For Artisan, the Java backend and Vue frontend are separated modules. 
Both modules have a 'build' task. Example of `frontend/MODULE.hcl`: 

```hcl
# file: frontend/MODULE.hcl
task "node:lts-alpine" "build" {
  script = "npm install && npm run build"
  ...
}

# file: backend/MODULE.hcl
task "go" "build" {
  script = "go build"  
}
``` 

When we run the `backend` build, the Artisan executes tasks for each module 
within an own docker container. Because there is dependency set, the frontend 
is build first and backend last.

```
artisan run //backend:build
```

Check the [simple demo repository](http://github.com/unravela/artisan-simple-demo)
for demonstration of small monorepo with one frontend and backend application.

## Installation
If you're **Linux** user, you can use the following command:

```bash
curl -sfL https://artisan.unravela.io/install.sh | sh
```

If you're **Mac OS** user with [Homebrew](https://brew.sh) installed, you can 
install Artisan with the command:

```bash
brew install unravela/tap/artisan
```

If you're **Windows** user, you can download [ZIP archive](https://github.com/unravela/artisan/releases/latest) directly or you can use [Scoop](https://scoop.sh/):

```
scoop bucket add unravela https://github.com/unravela/scoop-bucket
scoop install artisan
```
