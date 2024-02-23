// Defines the error messages used within FlyteAdmin categorized by common error codes.
package errors

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/wI2L/jsondiff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/flyteorg/flyte/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyte/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/flyteorg/flyte/flytestdlib/logger"
)

type FlyteAdminError interface {
	Error() string
	Code() codes.Code
	GRPCStatus() *status.Status
	WithDetails(details proto.Message) (FlyteAdminError, error)
	String() string
}
type flyteAdminErrorImpl struct {
	status *status.Status
}

func (e *flyteAdminErrorImpl) Error() string {
	return e.status.Message()
}

func (e *flyteAdminErrorImpl) Code() codes.Code {
	return e.status.Code()
}

func (e *flyteAdminErrorImpl) GRPCStatus() *status.Status {
	return e.status
}

func (e *flyteAdminErrorImpl) String() string {
	return fmt.Sprintf("status: %v", e.status)
}

// enclose the error in the format that grpc server expect from golang:
//
//	https://github.com/grpc/grpc-go/blob/master/status/status.go#L133
func (e *flyteAdminErrorImpl) WithDetails(details proto.Message) (FlyteAdminError, error) {
	s, err := e.status.WithDetails(details)
	if err != nil {
		return nil, err
	}
	return NewFlyteAdminErrorFromStatus(s), nil
}

func NewFlyteAdminErrorFromStatus(status *status.Status) FlyteAdminError {
	return &flyteAdminErrorImpl{
		status: status,
	}
}

func NewFlyteAdminError(code codes.Code, message string) FlyteAdminError {
	return &flyteAdminErrorImpl{
		status: status.New(code, message),
	}
}

func NewFlyteAdminErrorf(code codes.Code, format string, a ...interface{}) FlyteAdminError {
	return NewFlyteAdminError(code, fmt.Sprintf(format, a...))
}

func toStringSlice(errors []error) []string {
	errSlice := make([]string, len(errors))
	for idx, err := range errors {
		errSlice[idx] = err.Error()
	}
	return errSlice
}

func NewCollectedFlyteAdminError(code codes.Code, errors []error) FlyteAdminError {
	return NewFlyteAdminError(code, strings.Join(toStringSlice(errors), ", "))
}

func NewAlreadyInTerminalStateError(ctx context.Context, errorMsg string, curPhase string) FlyteAdminError {
	logger.Warn(ctx, errorMsg)
	alreadyInTerminalPhase := &admin.EventErrorAlreadyInTerminalState{CurrentPhase: curPhase}
	reason := &admin.EventFailureReason{
		Reason: &admin.EventFailureReason_AlreadyInTerminalState{AlreadyInTerminalState: alreadyInTerminalPhase},
	}
	statusErr, transformationErr := NewFlyteAdminError(codes.FailedPrecondition, errorMsg).WithDetails(reason)
	if transformationErr != nil {
		logger.Panicf(ctx, "Failed to wrap grpc status in type 'Error': %v", transformationErr)
		return NewFlyteAdminErrorf(codes.FailedPrecondition, errorMsg)
	}
	return statusErr
}

func NewIncompatibleClusterError(ctx context.Context, errorMsg, curCluster string) FlyteAdminError {
	statusErr, transformationErr := NewFlyteAdminError(codes.FailedPrecondition, errorMsg).WithDetails(&admin.EventFailureReason{
		Reason: &admin.EventFailureReason_IncompatibleCluster{
			IncompatibleCluster: &admin.EventErrorIncompatibleCluster{
				Cluster: curCluster,
			},
		},
	})
	if transformationErr != nil {
		logger.Panicf(ctx, "Failed to wrap grpc status in type 'Error': %v", transformationErr)
		return NewFlyteAdminErrorf(codes.FailedPrecondition, errorMsg)
	}
	return statusErr
}

func compareJsons(jsonArray1 jsondiff.Patch, jsonArray2 jsondiff.Patch) []string {
	results := []string{}
	map1 := make(map[string]jsondiff.Operation)
	for _, obj := range jsonArray1 {
		map1[obj.Path] = obj
	}

	for _, obj := range jsonArray2 {
		if val, ok := map1[obj.Path]; ok {
			result := fmt.Sprintf("%s: %s -> %s\t", obj.Path, obj.Value, val.Value)
			results = append(results, result)
		}
	}
	return results
}

func NewTaskExistsDifferentStructureError(ctx context.Context, request *admin.TaskCreateRequest, oldSpec *core.TaskTemplate, newSpec *core.TaskTemplate) FlyteAdminError {
	errorMsg := "task with different structure already exists:\n"
	// omit source code file object storage path
	diff, _ := jsondiff.Compare(oldSpec, newSpec, jsondiff.Ignores("/Target/Container/args/2"))
	rdiff, _ := jsondiff.Compare(newSpec, oldSpec, jsondiff.Ignores("/Target/Container/args/2"))
	rs := compareJsons(diff, rdiff)
	errorMsg += strings.Join(rs, "\n")

	statusErr, transformationErr := NewFlyteAdminError(codes.InvalidArgument, errorMsg).WithDetails(&admin.CreateTaskFailureReason{
		Reason: &admin.CreateTaskFailureReason_ExistsDifferentStructure{
			ExistsDifferentStructure: &admin.TaskErrorExistsDifferentStructure{
				Id: request.Id,
			},
		},
	})
	if transformationErr != nil {
		logger.Panicf(ctx, "Failed to wrap grpc status in type 'Error': %v", transformationErr)
		return NewFlyteAdminErrorf(codes.InvalidArgument, errorMsg)
	}
	return statusErr
}

func NewTaskExistsIdenticalStructureError(ctx context.Context, request *admin.TaskCreateRequest) FlyteAdminError {
	errorMsg := "task with identical structure already exists"
	statusErr, transformationErr := NewFlyteAdminError(codes.AlreadyExists, errorMsg).WithDetails(&admin.CreateTaskFailureReason{
		Reason: &admin.CreateTaskFailureReason_ExistsIdenticalStructure{
			ExistsIdenticalStructure: &admin.TaskErrorExistsIdenticalStructure{
				Id: request.Id,
			},
		},
	})
	if transformationErr != nil {
		logger.Panicf(ctx, "Failed to wrap grpc status in type 'Error': %v", transformationErr)
		return NewFlyteAdminErrorf(codes.AlreadyExists, errorMsg)
	}
	return statusErr
}

func NewWorkflowExistsDifferentStructureError(ctx context.Context, request *admin.WorkflowCreateRequest, oldTemplate *core.TaskTemplate, newTemplate *core.TaskTemplate) FlyteAdminError {
	errorMsg := "workflow with different structure already exists:\n"
	// omit source code file object storage path
	diff, _ := jsondiff.Compare(oldTemplate, newTemplate, jsondiff.Ignores("/Target/Container/args/2"))
	rdiff, _ := jsondiff.Compare(newTemplate, oldTemplate, jsondiff.Ignores("/Target/Container/args/2"))
	rs := compareJsons(diff, rdiff)
	errorMsg += strings.Join(rs, "\n")

	statusErr, transformationErr := NewFlyteAdminError(codes.InvalidArgument, errorMsg).WithDetails(&admin.CreateWorkflowFailureReason{
		Reason: &admin.CreateWorkflowFailureReason_ExistsDifferentStructure{
			ExistsDifferentStructure: &admin.WorkflowErrorExistsDifferentStructure{
				Id: request.Id,
			},
		},
	})
	if transformationErr != nil {
		logger.Panicf(ctx, "Failed to wrap grpc status in type 'Error': %v", transformationErr)
		return NewFlyteAdminErrorf(codes.InvalidArgument, errorMsg)
	}
	return statusErr
}

func NewWorkflowExistsIdenticalStructureError(ctx context.Context, request *admin.WorkflowCreateRequest) FlyteAdminError {
	errorMsg := "workflow with identical structure already exists"
	statusErr, transformationErr := NewFlyteAdminError(codes.AlreadyExists, errorMsg).WithDetails(&admin.CreateWorkflowFailureReason{
		Reason: &admin.CreateWorkflowFailureReason_ExistsIdenticalStructure{
			ExistsIdenticalStructure: &admin.WorkflowErrorExistsIdenticalStructure{
				Id: request.Id,
			},
		},
	})
	if transformationErr != nil {
		logger.Panicf(ctx, "Failed to wrap grpc status in type 'Error': %v", transformationErr)
		return NewFlyteAdminErrorf(codes.AlreadyExists, errorMsg)
	}
	return statusErr
}

func IsDoesNotExistError(err error) bool {
	adminError, ok := err.(FlyteAdminError)
	return ok && adminError.Code() == codes.NotFound
}
