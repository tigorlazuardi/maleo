# Minio

Maleo Minio is an implementation of Bucket wrapping around
[Minio SDK Go v7](https://min.io/docs/minio/linux/developers/go/minio-go.html).

## Installation

```bash
go get github.com/tigorlazuardi/maleo/bucket/maleominio/v2
```

## Usage

Usage is pretty straight forward. Create a client, and call `maleominio.Wrap()` around the client.

!!! warning "Requirements"

    Ensure your credential is allowed to execute `PutObject` Api on target bucket.

!!! note "Auto Bucket Create"

    If your credential support creating bucket, it will create the bucket if does not exist.

!!! info "Using on AWS"

    Minio client supports AWS S3 protocols. So you can also use this bucket implementation for AWS.

=== "Minio"

    ```go
    --8<-- "bucket/maleominio/v7/minio_example_test.go:example"
    ```

=== "AWS"

    ```go
    --8<-- "bucket/maleominio/v7/minio_example_test.go:aws"
    ```
