package helper

const (
	// Approval Status IDs (from references table)
	ApprovalStatusWaiting        = 20
	ApprovalStatusReject         = 21
	ApprovalStatusApprove        = 22
	ApprovalStatusPendingManager = 23

	// Task Status IDs
	TaskStatusToDo       = 1
	TaskStatusOnHold     = 2
	TaskStatusOnProgress = 3
	TaskStatusDone       = 4
	TaskStatusInReview   = 37
	TaskStatusRevision   = 39

	// File upload paths
	TaskFileUploadPath   = "file/tasks/"
	PermitFileUploadPath = "file/permits/"
)
