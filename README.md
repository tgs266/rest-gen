# rest-gen

rest-gen is a golang api generator, inspired by [conjure](https://github.com/palantir/conjure) and [OAS](https://www.openapis.org/)
## Motivation

Some might ask, why was this tool built? OAS exists and works pretty well, as do other tools. Unfortunately, some of these tools have some pretty annoying drawbacks, like Conjure doesn't support windows, and OAS generators for golang don't allow you to generate just types, or to split up your specs in a clean way. 
## Installation

```
go install github.com/tgs266/rest-gen@latest
```
## Running

To run, define your specs, and then run 
```
rest-gen -i INPUT_FOLDER -o OUTPUT_FOLDER
```

## Roadmap

In no order
* Add runtime package
* Better error handling