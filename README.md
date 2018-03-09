## dep-graphviz

The code is extracted from [dep](https://github.com/golang/dep).

The dep has supported for generated dependencies graph by call `dep status -dot`, but it only supported for the project that used dep for package manage and only for the whole project.

When reading some non-dep go projects or very large projects it can be unreadable.

This code can run under any folder and generated the dependencies graph without dep.

## Usage

```
# Make sure you have installed graphviz
go get -u github.com/jerry-tao/dep-graphviz/cmd/dg

# Run in current folder
dg | dot -Tpdf -O

# ignore stdlib
dg $GOPATH/src/MYPROJECT | dot -Tpdf -O

# include stdlib
dg -s $GOPATH/src/MYPROJECT | dot -Tpdf -O


```