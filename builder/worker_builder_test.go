package builder

import (
	"testing"

	"github.com/getumen/arachne"
	gomock "github.com/golang/mock/gomock"
)

func TestWorkerBuilder_BuildFail(t *testing.T) {
	builder := NewWorkerBuilder()
	_, err := builder.Build()
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
}

func TestWorkerBuilder_BuildSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpClientMock := arachne.NewMockHTTPClient(ctrl)
	loggerMock := arachne.NewMockLogger(ctrl)
	workerQueueMock := arachne.NewMockWorkerQueue(ctrl)

	builder := NewWorkerBuilder()
	builder.SetHTTPClient(httpClientMock)
	builder.SetLogger(loggerMock)
	builder.SetWorkerQueue(workerQueueMock)
	builder.SetSpider(func(*arachne.Response) ([]*arachne.Request, error) { return nil, nil })
	_, err := builder.Build()
	if err != nil {
		t.Fatalf("expected nil, but got error: %v", err)
	}
}
