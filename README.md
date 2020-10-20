# Artisan

[![Release](https://img.shields.io/github/release/unravela/artisan.svg?style=flat-square)](https://github.com/unravela/artisan/releases/latest)
[![Build](https://img.shields.io/github/workflow/status/unravela/artisan/build?style=flat-square)](https://github.com/unravela/artisan/actions?query=workflow%3Abuild)


Artisan helps you build complex heterogeneous repositories without the need to 
install any complicated tools of specific versions as pre-requirements. The main 
idea is: Run the same build everywhere.
 
## How it works
Let's have a repository with a backend application written as Java/Gradle 
and Vue.JS frontend. Usually, we need to install the correct versions of NPM, Java, 
and Gradle into our system. For Artisan, the backend and frontend are modules. Each 
module has `build` task executed in its own docker container. 

```hcl
# file: frontend/MODULE.hcl
task "node:lts-alpine" "build" {
  script = "npm install && npm run build"
  ...
}

# file: backend/MODULE.hcl
task "gradle:6.7.0-jdk11" "build" {
  script = "gradle build"
  deps = [ "//frontend:build" ]  
}
``` 

When we run the `backend` build, the Artisan executes tasks for each module 
within an own docker container. There is no need to have NPM or Java installed. 
Because there is dependency between backend and frontend, the frontend is build 
first and backend last.

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
