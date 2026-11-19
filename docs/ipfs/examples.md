# Client examples

All examples use ordinary S3 clients and a preassigned tenant prefix. The administrator must determine the caller's UUID from the internal tenant mapping; clients cannot choose it. Replace the example prefix with the authorized one.

## AWS CLI

```sh
aws --endpoint-url https://s3.example.com s3api put-object \
  --bucket shared --key users/UUID/photo.jpg --body photo.jpg \
  --metadata maxio-pin=true,maxio-access=public
aws --endpoint-url https://s3.example.com s3api head-object \
  --bucket shared --key users/UUID/photo.jpg
aws --endpoint-url https://s3.example.com s3api put-object-tagging \
  --bucket shared --key users/UUID/photo.jpg \
  --tagging 'TagSet=[{Key=maxio-pin,Value=false}]'
```

## rclone

Configure an S3 remote with MaxIO's endpoint and credentials, then copy into the authorized prefix:

```sh
rclone copy photo.jpg maxio:shared/users/UUID/
```

To request public pinning with rclone, supply its S3 metadata options or set the server's public automatic pinning policy. Check your installed rclone version for its current metadata flag syntax.

## Go AWS SDK v2

```go
_, err := client.PutObject(ctx, &s3.PutObjectInput{
	Bucket: aws.String("shared"),
	Key: aws.String("users/UUID/photo.jpg"),
	Body: file,
	Metadata: map[string]string{"maxio-pin": "true", "maxio-access": "public"},
})
if err != nil { return err }
head, err := client.HeadObject(ctx, &s3.HeadObjectInput{
	Bucket: aws.String("shared"), Key: aws.String("users/UUID/photo.jpg"),
})
if err != nil { return err }
_ = head.Metadata // Standard S3 user metadata remains available.
```

AWS SDK v2's `HeadObjectOutput` does not automatically map custom `x-maxio-*` headers. Use a response middleware or a raw signed request to read the IPFS status headers.
