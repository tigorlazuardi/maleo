package maleos3

import "fmt"

type URLBuilder interface {
	BuildURL(obj Object) string
}

type URLBuilderFunc func(obj Object) string

func (f URLBuilderFunc) BuildURL(obj Object) string {
	return f(obj)
}

func s3URLBuilder(bucketInEndpoint bool) URLBuilder {
	return URLBuilderFunc(func(obj Object) string {
		if bucketInEndpoint {
			return fmt.Sprintf("%s://%s/%s", obj.Proto, obj.Endpoint, obj.Key)
		}
		return fmt.Sprintf("%s://%s/%s/%s", obj.Proto, obj.Endpoint, obj.Bucket, obj.Key)
	})
}
