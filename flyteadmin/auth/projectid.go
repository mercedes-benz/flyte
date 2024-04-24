package auth

import (
	"github.com/flyteorg/flyte/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyte/flyteidl/gen/pb-go/flyteidl/service"
)

// Method to request type mapping is based on flyteidl/gen/pb-go/flyteidl/service/admin_grpc.pb.go
var methods = map[string]func(req interface{}) string{
	service.AdminService_CreateTask_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.TaskCreateRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_GetTask_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ObjectGetRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_ListTaskIds_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.NamedEntityIdentifierListRequest)
		return request.GetProject()
	},
	service.AdminService_ListTasks_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ResourceListRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_CreateWorkflow_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.WorkflowCreateRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_GetWorkflow_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ObjectGetRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_ListWorkflowIds_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.NamedEntityIdentifierListRequest)
		return request.GetProject()
	},
	service.AdminService_ListWorkflows_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ResourceListRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_CreateLaunchPlan_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.LaunchPlanCreateRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_GetLaunchPlan_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ObjectGetRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_GetActiveLaunchPlan_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ActiveLaunchPlanRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_ListActiveLaunchPlans_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ActiveLaunchPlanListRequest)
		return request.GetProject()
	},
	service.AdminService_ListLaunchPlanIds_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.NamedEntityIdentifierListRequest)
		return request.GetProject()
	},
	service.AdminService_ListLaunchPlans_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ResourceListRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_UpdateLaunchPlan_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.LaunchPlanUpdateRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_CreateExecution_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ExecutionCreateRequest)
		return request.GetProject()
	},
	service.AdminService_RelaunchExecution_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ExecutionRelaunchRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_RecoverExecution_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ExecutionRecoverRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_GetExecution_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.WorkflowExecutionGetRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_UpdateExecution_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ExecutionUpdateRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_GetExecutionData_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.WorkflowExecutionGetDataRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_ListExecutions_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ResourceListRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_TerminateExecution_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ExecutionTerminateRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_GetNodeExecution_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.NodeExecutionGetRequest)
		return request.GetId().GetExecutionId().GetProject()
	},
	service.AdminService_GetDynamicNodeWorkflow_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.GetDynamicNodeWorkflowRequest)
		return request.GetId().GetExecutionId().GetProject()
	},
	service.AdminService_ListNodeExecutions_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.NodeExecutionListRequest)
		return request.GetWorkflowExecutionId().GetProject()
	},
	service.AdminService_ListNodeExecutionsForTask_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.NodeExecutionForTaskListRequest)
		return request.GetTaskExecutionId().GetTaskId().GetProject()
	},
	service.AdminService_GetNodeExecutionData_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.NodeExecutionGetDataRequest)
		return request.GetId().GetExecutionId().GetProject()
	},
	service.AdminService_RegisterProject_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ProjectRegisterRequest)
		return request.GetProject().GetId()
	},
	service.AdminService_UpdateProject_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.Project)
		return request.GetId()
	},
	service.AdminService_ListProjects_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ProjectListRequest)
		return request.String()
	},
	service.AdminService_CreateWorkflowEvent_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.WorkflowExecutionEventRequest)
		return request.GetRequestId()
	},
	service.AdminService_CreateNodeEvent_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.NodeExecutionEventRequest)
		return request.GetRequestId()
	},
	service.AdminService_CreateTaskEvent_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.TaskExecutionEventRequest)
		return request.GetRequestId()
	},
	service.AdminService_GetTaskExecution_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.TaskExecutionGetRequest)
		return request.GetId().GetTaskId().GetProject()
	},
	service.AdminService_ListTaskExecutions_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.TaskExecutionListRequest)
		return request.GetNodeExecutionId().GetExecutionId().GetProject()
	},
	service.AdminService_GetTaskExecutionData_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.TaskExecutionGetDataRequest)
		return request.GetId().GetTaskId().GetProject()
	},
	service.AdminService_UpdateProjectDomainAttributes_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ProjectDomainAttributesUpdateRequest)
		return request.GetAttributes().GetProject()
	},
	service.AdminService_GetProjectDomainAttributes_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ProjectDomainAttributesGetRequest)
		return request.GetProject()
	},
	service.AdminService_DeleteProjectDomainAttributes_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ProjectDomainAttributesDeleteRequest)
		return request.GetProject()
	},
	service.AdminService_UpdateProjectAttributes_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ProjectAttributesUpdateRequest)
		return request.GetAttributes().GetProject()
	},
	service.AdminService_GetProjectAttributes_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ProjectAttributesGetRequest)
		return request.GetProject()
	},
	service.AdminService_DeleteProjectAttributes_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ProjectAttributesDeleteRequest)
		return request.GetProject()
	},
	service.AdminService_UpdateWorkflowAttributes_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.WorkflowAttributesUpdateRequest)
		return request.GetAttributes().GetProject()
	},
	service.AdminService_GetWorkflowAttributes_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.WorkflowAttributesGetRequest)
		return request.GetProject()
	},
	service.AdminService_DeleteWorkflowAttributes_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.WorkflowAttributesDeleteRequest)
		return request.GetProject()
	},
	// service.AdminService_ListMatchableAttributes_FullMethodName: func(req interface{}) string {
	// 	request, _ := req.(*admin.ListMatchableAttributesRequest)
	// 	return request.GetProject()
	// },
	service.AdminService_ListNamedEntities_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.NamedEntityListRequest)
		return request.GetProject()
	},
	service.AdminService_GetNamedEntity_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.NamedEntityGetRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_UpdateNamedEntity_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.NamedEntityUpdateRequest)
		return request.GetId().GetProject()
	},
	// service.AdminService_GetVersion_FullMethodName: func(req interface{}) string {
	// 	request, _ := req.(*admin.GetVersionRequest)
	// 	return request.GetProject()
	// },
	service.AdminService_GetDescriptionEntity_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.ObjectGetRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_ListDescriptionEntities_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.DescriptionEntityListRequest)
		return request.GetId().GetProject()
	},
	service.AdminService_GetExecutionMetrics_FullMethodName: func(req interface{}) string {
		request, _ := req.(*admin.WorkflowExecutionGetMetricsRequest)
		return request.GetId().GetProject()
	},
}

func inferProjectIDFromAdminRequest(fullMethod string, req interface{}) string {
	method := methods[fullMethod]
	if method != nil {
		return method(req)
	}
	return ""
}
