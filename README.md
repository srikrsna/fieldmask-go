# fieldmask-go

Generate helpers to construct compile time safe fieldmasks in go.

Install the plugin using,

```bash
go get github.com/srikrsna/fieldmask-go/cmd/protoc-gen-fieldmask-go@latest
```

Use it:

```yaml
version: v1
plugins:
  - name: go
    out: gen
    opt: paths=source_relative
  - name: fieldmask-go
    out: gen
    opt: paths=source_relative
```

For protobuf messages,

```proto
message Entity {
    string id = 1;
    SubMessage sub = 2;
}

message SubMessage {
    string id = 1;
}
```

This should now generate code that can be used,

```go


fm, err := fieldmaskpb.New(
    &pb.Entity{},
    string(pbfieldmask.EntityMask.Id()),
    pbfieldmask.EntityMask.Sub().Id(),
)

```
