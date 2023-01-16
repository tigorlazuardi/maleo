package maleos3

// Object is a simple metadata of an object in S3.
type Object struct {
	Endpoint string
	Bucket   string
	Key      string
	Region   string
	Proto    string
}

// NewObject is a constructor for Object.
func NewObject(proto, endpoint, region, bucket, key string) Object {
	return Object{
		Endpoint: endpoint,
		Bucket:   bucket,
		Key:      key,
		Region:   region,
		Proto:    proto,
	}
}
