---
hide:
    - toc
---

# Minio (V7)

Maleo Minio is an implementation of Bucket wrapping around
[Minio SDK Go v7](https://min.io/docs/minio/linux/developers/go/minio-go.html).

## Installation

```bash
go get github.com/tigorlazuardi/maleo/bucket/maleominio-v7
```

## Usage

Usage is pretty straight forward. Create a client, and call `maleominio.Wrap()` around the client.

You might want to consult [Minio's docs](https://min.io/docs/minio/linux/developers/go/minio-go.html#id4) regarding the
setup.

!!! warning "Requirements"

    Ensure your credential has permission to execute `PutObject` API on target bucket.

!!! note "Auto Bucket Create"

    If your credential support creating bucket, it will create the bucket if does not exist.

=== "Minio"

    ```go
    --8<-- "bucket/maleominio-v7/minio_example_test.go:example"
    ```

=== "AWS"

    ```go
    --8<-- "bucket/maleominio-v7/minio_example_test.go:aws"
    ```
