package maleos3

import "testing"

func Test_detectRegion(t *testing.T) {
	type args struct {
		endpoint string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "detect region from endpoint",
			args: args{
				endpoint: "my-bucket.s3.us-east-1.amazonaws.com",
			},
			want: "us-east-1",
		},
		{
			name: "detect region from endpoint (no bucket)",
			args: args{
				endpoint: "s3.ap-southeast-1.amazonaws.com",
			},
			want: "ap-southeast-1",
		},
		{
			name: "ignore endpoint too short",
			args: args{
				endpoint: "amazonaws.com",
			},
		},
		{
			name: "ignore endpoint with bucket but no region",
			args: args{
				endpoint: "my-bucket.s3.amazonaws.com",
			},
		},
		{
			name: "ignore non AWS FQDN",
			args: args{
				endpoint: "my-bucket.s3.my-region.my-domain.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectRegion(tt.args.endpoint); got != tt.want {
				t.Errorf("detectRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_detectBucket(t *testing.T) {
	type args struct {
		endpoint string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "detect bucket from endpoint",
			args: args{
				endpoint: "my-bucket.s3.us-east-1.amazonaws.com",
			},
			want: "my-bucket",
		},
		{
			name: "ignore endpoint without region (invalid endpoint from AWS perspective)",
			args: args{
				endpoint: "my-bucket.s3.amazonaws.com",
			},
		},
		{
			name: "ignore non AWS FQDN",
			args: args{
				endpoint: "my-bucket.s3.my-region.my-domain.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectBucket(tt.args.endpoint); got != tt.want {
				t.Errorf("detectBucket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isAws(t *testing.T) {
	type args struct {
		endpoint string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "detect AWS FQDN",
			args: args{
				endpoint: "my-bucket.s3.us-east-1.amazonaws.com",
			},
			want: true,
		},
		{
			name: "false on non amazonaws.com",
			args: args{
				endpoint: "my-bucket.s3.my-region.my-domain.com",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAws(tt.args.endpoint); got != tt.want {
				t.Errorf("isAws() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewS3BucketErrors(t *testing.T) {
	type args struct {
		endpoint string
		opts     []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "endpoint empty - no bucket",
			args: args{
				endpoint: "",
			},
			wantErr: true,
		},
		{
			name: "endpoint empty - no region",
			args: args{
				endpoint: "",
				opts:     []Option{WithBucket("foo")},
			},
			wantErr: true,
		},
		{
			name: "no bucket in endpoint without WithBucket resolves as error",
			args: args{
				endpoint: "s3.us-east-1.amazonaws.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewS3Bucket(tt.args.endpoint, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewS3Bucket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
