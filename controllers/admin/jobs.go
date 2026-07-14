package admin

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mortenvistisen/internal/inertia"
	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/queue"
	"mortenvistisen/router"
	"mortenvistisen/router/cookies"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v5"
	"github.com/riverqueue/river/rivertype"
)

const jobPageSize int64 = 25

type Jobs struct {
	db        storage.Pool
	processor queue.Processor
}

func NewJobs(db storage.Pool, processor queue.Processor) Jobs {
	return Jobs{db: db, processor: processor}
}

func (j Jobs) RegisterRoutes(r *router.Router) error {
	registrations := []echo.Route{
		{
			Method:  http.MethodGet,
			Path:    routes.AdminJobIndex.Path(),
			Name:    routes.AdminJobIndex.Name(),
			Handler: j.Index,
		},
		{
			Method:  http.MethodGet,
			Path:    routes.AdminJobShow.Path(),
			Name:    routes.AdminJobShow.Name(),
			Handler: j.Show,
		},
		{
			Method:  http.MethodPost,
			Path:    routes.AdminJobRun.Path(),
			Name:    routes.AdminJobRun.Name(),
			Handler: j.Run,
		},
		{
			Method:  http.MethodPost,
			Path:    routes.AdminJobRetry.Path(),
			Name:    routes.AdminJobRetry.Name(),
			Handler: j.Retry,
		},
		{
			Method:  http.MethodPost,
			Path:    routes.AdminJobCancel.Path(),
			Name:    routes.AdminJobCancel.Name(),
			Handler: j.Cancel,
		},
		{
			Method:  http.MethodDelete,
			Path:    routes.AdminJobDestroy.Path(),
			Name:    routes.AdminJobDestroy.Name(),
			Handler: j.Destroy,
		},
	}

	var errs []error
	for _, registration := range registrations {
		if _, err := r.AddRoute(registration); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

type JobSummaryData struct {
	ID          int64
	State       string
	Attempt     int16
	MaxAttempts int16
	AttemptedAt *time.Time
	CreatedAt   time.Time
	FinalizedAt *time.Time
	ScheduledAt time.Time
	Priority    int16
	Kind        string
	Queue       string
}

func newJobSummaryData(entity models.JobSummary) JobSummaryData {
	var attemptedAt *time.Time
	if entity.AttemptedAt.Valid {
		attemptedAt = &entity.AttemptedAt.Time
	}
	var finalizedAt *time.Time
	if entity.FinalizedAt.Valid {
		finalizedAt = &entity.FinalizedAt.Time
	}

	return JobSummaryData{
		ID:          entity.ID,
		State:       entity.State,
		Attempt:     entity.Attempt,
		MaxAttempts: entity.MaxAttempts,
		AttemptedAt: attemptedAt,
		CreatedAt:   entity.CreatedAt,
		FinalizedAt: finalizedAt,
		ScheduledAt: entity.ScheduledAt,
		Priority:    entity.Priority,
		Kind:        entity.Kind,
		Queue:       entity.Queue,
	}
}

func newJobSummaryDataList(entities []models.JobSummary) []JobSummaryData {
	items := make([]JobSummaryData, 0, len(entities))
	for _, entity := range entities {
		items = append(items, newJobSummaryData(entity))
	}
	return items
}

type JobData struct {
	ID           int64
	State        string
	Attempt      int
	MaxAttempts  int
	AttemptedAt  *time.Time
	AttemptedBy  []string
	CreatedAt    time.Time
	FinalizedAt  *time.Time
	ScheduledAt  time.Time
	Priority     int
	Args         json.RawMessage
	Errors       []rivertype.AttemptError
	Kind         string
	Metadata     json.RawMessage
	Queue        string
	Tags         []string
	UniqueKey    string
	UniqueStates []rivertype.JobState
	CanRun       bool
	CanRestart   bool
	CanCancel    bool
	CanDelete    bool
}

func newJobData(job *rivertype.JobRow) JobData {
	attemptedBy := job.AttemptedBy
	if attemptedBy == nil {
		attemptedBy = []string{}
	}
	errorsList := job.Errors
	if errorsList == nil {
		errorsList = []rivertype.AttemptError{}
	}
	tags := job.Tags
	if tags == nil {
		tags = []string{}
	}
	uniqueStates := job.UniqueStates
	if uniqueStates == nil {
		uniqueStates = []rivertype.JobState{}
	}

	data := JobData{
		ID:           job.ID,
		State:        string(job.State),
		Attempt:      job.Attempt,
		MaxAttempts:  job.MaxAttempts,
		AttemptedAt:  job.AttemptedAt,
		AttemptedBy:  attemptedBy,
		CreatedAt:    job.CreatedAt,
		FinalizedAt:  job.FinalizedAt,
		ScheduledAt:  job.ScheduledAt,
		Priority:     job.Priority,
		Args:         json.RawMessage(job.EncodedArgs),
		Errors:       errorsList,
		Kind:         job.Kind,
		Metadata:     json.RawMessage(job.Metadata),
		Queue:        job.Queue,
		Tags:         tags,
		UniqueKey:    hex.EncodeToString(job.UniqueKey),
		UniqueStates: uniqueStates,
		CanDelete:    job.State != rivertype.JobStateRunning,
	}

	switch job.State {
	case rivertype.JobStateScheduled, rivertype.JobStateRetryable, rivertype.JobStatePending:
		data.CanRun = true
	case rivertype.JobStateCancelled, rivertype.JobStateCompleted, rivertype.JobStateDiscarded:
		data.CanRestart = true
	case rivertype.JobStateAvailable, rivertype.JobStateRunning:
	}

	switch job.State {
	case rivertype.JobStateAvailable, rivertype.JobStatePending, rivertype.JobStateRetryable,
		rivertype.JobStateRunning, rivertype.JobStateScheduled:
		data.CanCancel = true
	case rivertype.JobStateCancelled, rivertype.JobStateCompleted, rivertype.JobStateDiscarded:
	}

	return data
}

type JobPaginationData struct {
	Page       int64
	PageSize   int64
	TotalCount int64
	TotalPages int64
}

type JobFiltersData struct {
	State  string
	Search string
}

func validJobState(value string) bool {
	for _, state := range rivertype.JobStates() {
		if value == string(state) {
			return true
		}
	}
	return false
}

func (j Jobs) Index(etx *echo.Context) error {
	page := int64(1)
	if parsed, err := strconv.ParseInt(etx.QueryParam("page"), 10, 64); err == nil && parsed > 0 {
		page = parsed
	}

	filter := models.JobFilter{Search: strings.TrimSpace(etx.QueryParam("search"))}
	if state := etx.QueryParam("state"); validJobState(state) {
		filter.State = state
	}

	jobList, err := models.Job.Paginate(
		etx.Request().Context(),
		j.db.Executor(),
		filter,
		page,
		jobPageSize,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	stats, err := models.Job.Stats(etx.Request().Context(), j.db.Executor())
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Job/Index", inertia.Props{
		"items": newJobSummaryDataList(jobList.Jobs),
		"pagination": JobPaginationData{
			Page:       jobList.Page,
			PageSize:   jobList.PageSize,
			TotalCount: jobList.TotalCount,
			TotalPages: jobList.TotalPages,
		},
		"filters": JobFiltersData{State: filter.State, Search: filter.Search},
		"stats":   stats,
	})
}

func (j Jobs) Show(etx *echo.Context) error {
	jobID, err := jobIDParam(etx)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	job, err := j.processor.Job(etx.Request().Context(), jobID)
	if errors.Is(err, rivertype.ErrNotFound) {
		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Job/Show", inertia.Props{"item": newJobData(job)})
}

func jobIDParam(etx *echo.Context) (int64, error) {
	return strconv.ParseInt(etx.Param("id"), 10, 64)
}

func jobActionError(etx *echo.Context, jobID int64, action string, err error) error {
	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashError,
		fmt.Sprintf("Could not %s job: %v", action, err),
	); flashErr != nil {
		return flashErr
	}
	return inertia.Redirect(etx, routes.AdminJobShow.URL(jobID))
}

func jobActionSuccess(etx *echo.Context, jobID int64, message string) error {
	if err := cookies.AddFlash(etx, cookies.FlashSuccess, message); err != nil {
		return err
	}
	return inertia.Redirect(etx, routes.AdminJobShow.URL(jobID))
}

func (j Jobs) Run(etx *echo.Context) error {
	jobID, err := jobIDParam(etx)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	if _, err := j.processor.RunJobNow(etx.Request().Context(), jobID); err != nil {
		return jobActionError(etx, jobID, "run", err)
	}
	return jobActionSuccess(etx, jobID, "Job made available to run now")
}

func (j Jobs) Retry(etx *echo.Context) error {
	jobID, err := jobIDParam(etx)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	if _, err := j.processor.RestartJob(etx.Request().Context(), jobID); err != nil {
		return jobActionError(etx, jobID, "restart", err)
	}
	return jobActionSuccess(etx, jobID, "Job restarted and made available")
}

func (j Jobs) Cancel(etx *echo.Context) error {
	jobID, err := jobIDParam(etx)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	if _, err := j.processor.CancelJob(etx.Request().Context(), jobID); err != nil {
		return jobActionError(etx, jobID, "cancel", err)
	}
	return jobActionSuccess(etx, jobID, "Job cancellation requested")
}

func (j Jobs) Destroy(etx *echo.Context) error {
	jobID, err := jobIDParam(etx)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	if _, err := j.processor.DeleteJob(etx.Request().Context(), jobID); err != nil {
		return jobActionError(etx, jobID, "delete", err)
	}
	if err := cookies.AddFlash(etx, cookies.FlashSuccess, "Job permanently deleted"); err != nil {
		return err
	}
	return inertia.Redirect(etx, routes.AdminJobIndex.URL())
}
