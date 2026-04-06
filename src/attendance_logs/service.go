package attendanceLog

import (
	"context"
	"errors"

	"gin-auth-supabase/src/db"
)

type Service struct {
	q *db.Queries
}

func NewService(q *db.Queries) *Service {
	return &Service{q: q}
}

func (s *Service) Add(ctx context.Context, req AttendanceLogAdd) (*db.AttendanceLog, error) {

	attendanceLog, err := s.q.CreateAttendanceLog(ctx, db.CreateAttendanceLogParams{
		AssemblyPointID: req.AssemblyPointID,
		PersonnelCount:  req.PersonnelCount,
	})
	return &attendanceLog, err
}

func (s *Service) Request(ctx context.Context, req AttendanceLogRequest) (*[]db.AttendanceLog, error) {
	attendanceLogs, err := s.q.GetAttendanceLogByAssemblyPoint(ctx, req.AssemblyPointID)
	if err != nil {
		return nil, errors.New("Assembly Point not found")
	}

	return &attendanceLogs, nil
}
