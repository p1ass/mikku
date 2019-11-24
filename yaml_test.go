package mikku

import (
	"errors"
	"testing"
)

func Test_getCurrentTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		yamlStr   string
		imageName string
		want      string
		wantErr   error
	}{
		{
			name: "contains image",
			yamlStr: `
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: test-deployment
  labels:
    app: test-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      containers:
        - name: test-container
          image: asia.gcr.io/test/test-container:v1.0.0
          ports:
            - containerPort: 8080
`,
			imageName: "asia.gcr.io/test/test-container",
			want:      "v1.0.0",
			wantErr:   nil,
		},
		{
			name: "not contain image",
			yamlStr: `
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: test-deployment
  labels:
    app: test-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      containers:
        - name: test-container
          image: asia.gcr.io/test/test-container:v1.0.0
          ports:
            - containerPort: 8080
`,
			imageName: "asia.gcr.io/test/not-exist-container",
			want:      "",
			wantErr:   errImageNotFoundInYAML,
		},
		{
			name:      "invalid yaml",
			yamlStr:   ``,
			imageName: "asia.gcr.io/test/not-exist-container",
			want:      "",
			wantErr:   errImageNotFoundInYAML,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCurrentTag(tt.yamlStr, tt.imageName)
			if (tt.wantErr == nil && err != nil) || (tt.wantErr != nil && !errors.As(err, &tt.wantErr)) {
				t.Errorf("getCurrentTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getCurrentTag() = %v, wantContent %v", got, tt.want)
			}
		})
	}
}
