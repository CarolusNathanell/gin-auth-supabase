package attendanceLog

type AttendanceLogAdd struct {
	AssemblyPointID int32 `json:"assembly_point_id" binding:"required"`
	PersonnelCount  int32 `json:"personnel_count" binding:"required"`
}

type AttendanceLogRequest struct {
	AssemblyPointID int32 `json:"assembly_point_" binding:"required"`
}

type AttendanceLogResponse struct {
	AttendanceLog any `json:"attendance_log"`
}
