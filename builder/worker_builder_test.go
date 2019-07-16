package builder

import (
	"testing"

	"github.com/getumen/lucy"
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

	httpClientMock := lucy.NewMockHTTPClient(ctrl)
	loggerMock := lucy.NewMockLogger(ctrl)
	workerQueueMock := lucy.NewMockWorkerQueue(ctrl)

	builder := NewWorkerBuilder()
	builder.SetHTTPClient(httpClientMock)
	builder.SetLogger(loggerMock)
	builder.SetWorkerQueue(workerQueueMock)
	builder.SetSpider(func(*lucy.Response) ([]*lucy.Request, error) { return nil, nil })
	_, err := builder.Build()
	if err != nil {
		t.Fatalf("expected nil, but got error: %v", err)
	}
}
