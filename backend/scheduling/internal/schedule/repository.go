package schedule

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	h "github.com/maksmelnyk/scheduling/internal/database"
	e "github.com/maksmelnyk/scheduling/internal/database/entities"
)

type ScheduleRepo struct {
	db *sqlx.DB
}

func NewScheduleRepository(db *sqlx.DB) *ScheduleRepo {
	return &ScheduleRepo{db: db}
}

// GetWorkingPeriods retrieves working periods for a specific user within a date range
func (r *ScheduleRepo) GetWorkingPeriods(ctx context.Context, userId uuid.UUID, fromDate, toDate time.Time) ([]*e.WorkingPeriod, error) {
	const query = `
        SELECT id, user_id, date, start_time, end_time, created_at, updated_at
        FROM working_period
        WHERE user_id = $1 AND date BETWEEN $2 AND $3
    `
	return h.FetchMultiple[e.WorkingPeriod](ctx, r.db, query, userId, fromDate, toDate)
}

// GetScheduledEvents retrieves scheduled events for working periods
func (r *ScheduleRepo) GetScheduledEvents(ctx context.Context, workingPeriodIds []int64) ([]*e.ScheduledEvent, error) {
	const query = `
        SELECT id, session_id, working_period_id, user_id, start_time, end_time, max_participants, created_at, updated_at
        FROM scheduled_event
        WHERE working_period_id = ANY($1)
    `
	return h.FetchMultiple[e.ScheduledEvent](ctx, r.db, query, pq.Array(workingPeriodIds))
}

// GetBookings retrieves bookings for working period
func (r *ScheduleRepo) GetBookings(ctx context.Context, workingPeriodIds []int64) ([]*e.Booking, error) {
	const query = `
        SELECT id, teacher_id, student_id, session_Id, scheduled_event_id, working_period_id, start_time, end_time, status, created_at, updated_at
        FROM booking
        WHERE working_period_id = ANY($1)
    `
	return h.FetchMultiple[e.Booking](ctx, r.db, query, pq.Array(workingPeriodIds))
}

// GetWorkingPeriodById retrieves a single working period by its ID
func (r *ScheduleRepo) GetWorkingPeriodById(ctx context.Context, userId uuid.UUID, id int64) (*e.WorkingPeriod, error) {
	const query = `
        SELECT id, user_id, date, start_time, end_time, created_at, updated_at
        FROM working_period
        WHERE user_id = $1 AND id = $2
    `
	return h.FetchSingle[e.WorkingPeriod](ctx, r.db, query, userId, id)
}

// GetScheduledEventById retrieves a single scheduled event by its ID
func (r *ScheduleRepo) GetScheduledEventById(ctx context.Context, userId uuid.UUID, id int64) (*e.ScheduledEvent, error) {
	const query = `
		SELECT id, session_id, date, start_time, end_time, min_participants, max_participants, created_at, updated_at
		FROM scheduled_event
		WHERE session_id = $1 AND id = $2
	`
	return h.FetchSingle[e.ScheduledEvent](ctx, r.db, query, userId, id)
}

// HasLinkedEvents checks if a booking or scheduled event exists on a specific working period
func (r *ScheduleRepo) HasLinkedEvents(ctx context.Context, workingPeriodId int64) (bool, error) {
	const query = `
		SELECT 
			EXISTS (SELECT 1 FROM booking WHERE working_period_id = $1) OR 
			EXISTS (SELECT 1 FROM scheduled_event WHERE working_period_id = $1);
	`
	return h.CheckExists(ctx, r.db, query, workingPeriodId)
}

// AddWorkingPeriod adds a new working period
func (r *ScheduleRepo) AddWorkingPeriod(ctx context.Context, wd *e.WorkingPeriod) error {
	const query = `
        INSERT INTO working_period (user_id, date, start_time, end_time, created_at, updated_at)
        VALUES (:user_id, :date, :start_time, :end_time, :created_at, :updated_at)
    `
	return h.ExecNamedQuery(ctx, r.db, query, wd)
}

// UpdateWorkingPeriod updates an existing working period
func (r *ScheduleRepo) UpdateWorkingPeriod(ctx context.Context, wd *e.WorkingPeriod) error {
	const query = `
        UPDATE working_period
        SET date = :date, start_time = :start_time, end_time = :end_time, updated_at = :updated_at
        WHERE id = :id
    `
	return h.ExecNamedQuery(ctx, r.db, query, wd)
}

// AddScheduledEvent adds a new scheduled event
func (r *ScheduleRepo) AddScheduledEvent(ctx context.Context, se *e.ScheduledEvent) error {
	const query = `
		INSERT INTO scheduled_event (session_id, working_period_id, user_id, start_time, end_time, max_participants, created_at, updated_at)
		VALUES (:session_id, :working_period_id, :user_id, :start_time, :end_time, :max_participants, :created_at, :updated_at)
	`
	return h.ExecNamedQuery(ctx, r.db, query, se)
}

// DeleteWorkingPeriod deletes a working period by its ID
func (r *ScheduleRepo) DeleteWorkingPeriod(ctx context.Context, userId uuid.UUID, id int64) error {
	const query = `
        DELETE FROM working_period
        WHERE user_id = $1 AND id = $2
    `
	return h.ExecQuery(ctx, r.db, query, userId, id)
}

// DeleteScheduledEvent deletes a scheduled event by its ID
func (r *ScheduleRepo) DeleteScheduledEvent(ctx context.Context, userId uuid.UUID, id int64) error {
	const query = `
		DELETE FROM scheduled_event
		WHERE session_id = $1 AND id = $2
	`
	return h.ExecQuery(ctx, r.db, query, userId, id)
}
