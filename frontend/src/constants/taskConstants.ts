// Task Constants - matching backend helper/constants.go

// Approval Status IDs (from references table)
export const APPROVAL_STATUS = {
  WAITING: 20,
  REJECT: 21,
  APPROVE: 22,
  PENDING_MANAGER: 23,
} as const;

// Task Status IDs
export const TASK_STATUS = {
  TO_DO: 1,
  ON_HOLD: 2,
  ON_PROGRESS: 3,
  DONE: 4,
  IN_REVIEW: 37,
  REVISION: 39,
} as const;

// Task Type IDs
export const TASK_TYPE = {
  MAINTENANCE: 24,
  DEVELOPMENT: 25,
} as const;

// File Type IDs
export const FILE_TYPE = {
  CREATE: 30,
  BEFORE: 31,
  REVISION: 38,
} as const;
