package tests

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"go.temporal.io/server/common/convert"
	"go.temporal.io/server/common/persistence/sql/sqlplugin"
	"go.temporal.io/server/common/primitives"
	"go.temporal.io/server/common/shuffle"
)

type (
	historyExecutionActivitySuite struct {
		suite.Suite
		*require.Assertions

		store sqlplugin.HistoryExecutionActivity
	}
)

const (
	testHistoryExecutionActivityEncoding = "random encoding"
)

var (
	testHistoryExecutionActivityData = []byte("random history execution activity data")
)

func newHistoryExecutionActivitySuite(
	t *testing.T,
	store sqlplugin.HistoryExecutionActivity,
) *historyExecutionActivitySuite {
	return &historyExecutionActivitySuite{
		Assertions: require.New(t),
		store:      store,
	}
}

func (s *historyExecutionActivitySuite) SetupSuite() {

}

func (s *historyExecutionActivitySuite) TearDownSuite() {

}

func (s *historyExecutionActivitySuite) SetupTest() {
	s.Assertions = require.New(s.T())
}

func (s *historyExecutionActivitySuite) TearDownTest() {

}

func (s *historyExecutionActivitySuite) TestReplace_Single() {
	shardID := rand.Int31()
	namespaceID := primitives.NewUUID()
	workflowID := shuffle.String(testHistoryExecutionWorkflowID)
	runID := primitives.NewUUID()
	scheduleID := rand.Int63()

	activity := s.newRandomExecutionActivityRow(shardID, namespaceID, workflowID, runID, scheduleID)
	result, err := s.store.ReplaceIntoActivityInfoMaps([]sqlplugin.ActivityInfoMapsRow{activity})
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))
}

func (s *historyExecutionActivitySuite) TestReplace_Multiple() {
	shardID := rand.Int31()
	namespaceID := primitives.NewUUID()
	workflowID := shuffle.String(testHistoryExecutionWorkflowID)
	runID := primitives.NewUUID()

	activity1 := s.newRandomExecutionActivityRow(shardID, namespaceID, workflowID, runID, rand.Int63())
	activity2 := s.newRandomExecutionActivityRow(shardID, namespaceID, workflowID, runID, rand.Int63())
	result, err := s.store.ReplaceIntoActivityInfoMaps([]sqlplugin.ActivityInfoMapsRow{activity1, activity2})
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(2, int(rowsAffected))
}

func (s *historyExecutionActivitySuite) TestReplaceSelect_Single() {
	shardID := rand.Int31()
	namespaceID := primitives.NewUUID()
	workflowID := shuffle.String(testHistoryExecutionWorkflowID)
	runID := primitives.NewUUID()

	activity := s.newRandomExecutionActivityRow(shardID, namespaceID, workflowID, runID, rand.Int63())
	result, err := s.store.ReplaceIntoActivityInfoMaps([]sqlplugin.ActivityInfoMapsRow{activity})
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	filter := &sqlplugin.ActivityInfoMapsFilter{
		ShardID:     shardID,
		NamespaceID: namespaceID,
		WorkflowID:  workflowID,
		RunID:       runID,
	}
	rows, err := s.store.SelectFromActivityInfoMaps(filter)
	s.NoError(err)
	rowMap := map[int64]sqlplugin.ActivityInfoMapsRow{}
	for _, activity := range rows {
		rowMap[activity.ScheduleID] = activity
	}
	s.Equal(map[int64]sqlplugin.ActivityInfoMapsRow{
		activity.ScheduleID: activity,
	}, rowMap)
}

func (s *historyExecutionActivitySuite) TestReplaceSelect_Multiple() {
	numActivities := 20

	shardID := rand.Int31()
	namespaceID := primitives.NewUUID()
	workflowID := shuffle.String(testHistoryExecutionWorkflowID)
	runID := primitives.NewUUID()

	var activities []sqlplugin.ActivityInfoMapsRow
	for i := 0; i < numActivities; i++ {
		activity := s.newRandomExecutionActivityRow(shardID, namespaceID, workflowID, runID, rand.Int63())
		activities = append(activities, activity)
	}
	result, err := s.store.ReplaceIntoActivityInfoMaps(activities)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(numActivities, int(rowsAffected))

	filter := &sqlplugin.ActivityInfoMapsFilter{
		ShardID:     shardID,
		NamespaceID: namespaceID,
		WorkflowID:  workflowID,
		RunID:       runID,
	}
	rows, err := s.store.SelectFromActivityInfoMaps(filter)
	s.NoError(err)
	activityMap := map[int64]sqlplugin.ActivityInfoMapsRow{}
	for _, activity := range activities {
		activityMap[activity.ScheduleID] = activity
	}
	rowMap := map[int64]sqlplugin.ActivityInfoMapsRow{}
	for _, activity := range rows {
		rowMap[activity.ScheduleID] = activity
	}
	s.Equal(activityMap, rowMap)
}

func (s *historyExecutionActivitySuite) TestDeleteSelect_Single() {
	shardID := rand.Int31()
	namespaceID := primitives.NewUUID()
	workflowID := shuffle.String(testHistoryExecutionWorkflowID)
	runID := primitives.NewUUID()
	scheduleID := rand.Int63()

	filter := &sqlplugin.ActivityInfoMapsFilter{
		ShardID:     shardID,
		NamespaceID: namespaceID,
		WorkflowID:  workflowID,
		RunID:       runID,
		ScheduleID:  convert.Int64Ptr(scheduleID),
	}
	result, err := s.store.DeleteFromActivityInfoMaps(filter)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(0, int(rowsAffected))

	rows, err := s.store.SelectFromActivityInfoMaps(filter)
	s.NoError(err)
	s.Equal([]sqlplugin.ActivityInfoMapsRow(nil), rows)
}

func (s *historyExecutionActivitySuite) TestDeleteSelect_Multiple() {
	shardID := rand.Int31()
	namespaceID := primitives.NewUUID()
	workflowID := shuffle.String(testHistoryExecutionWorkflowID)
	runID := primitives.NewUUID()

	filter := &sqlplugin.ActivityInfoMapsFilter{
		ShardID:     shardID,
		NamespaceID: namespaceID,
		WorkflowID:  workflowID,
		RunID:       runID,
		ScheduleID:  nil,
	}
	result, err := s.store.DeleteFromActivityInfoMaps(filter)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(0, int(rowsAffected))

	rows, err := s.store.SelectFromActivityInfoMaps(filter)
	s.NoError(err)
	s.Equal([]sqlplugin.ActivityInfoMapsRow(nil), rows)
}

func (s *historyExecutionActivitySuite) TestReplaceDeleteSelect_Single() {
	shardID := rand.Int31()
	namespaceID := primitives.NewUUID()
	workflowID := shuffle.String(testHistoryExecutionWorkflowID)
	runID := primitives.NewUUID()
	scheduleID := rand.Int63()

	activity := s.newRandomExecutionActivityRow(shardID, namespaceID, workflowID, runID, scheduleID)
	result, err := s.store.ReplaceIntoActivityInfoMaps([]sqlplugin.ActivityInfoMapsRow{activity})
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	filter := &sqlplugin.ActivityInfoMapsFilter{
		ShardID:     shardID,
		NamespaceID: namespaceID,
		WorkflowID:  workflowID,
		RunID:       runID,
		ScheduleID:  convert.Int64Ptr(scheduleID),
	}
	result, err = s.store.DeleteFromActivityInfoMaps(filter)
	s.NoError(err)
	rowsAffected, err = result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	rows, err := s.store.SelectFromActivityInfoMaps(filter)
	s.NoError(err)
	s.Equal([]sqlplugin.ActivityInfoMapsRow(nil), rows)
}

func (s *historyExecutionActivitySuite) TestReplaceDeleteSelect_Multiple() {
	numActivities := 20

	shardID := rand.Int31()
	namespaceID := primitives.NewUUID()
	workflowID := shuffle.String(testHistoryExecutionWorkflowID)
	runID := primitives.NewUUID()

	var activities []sqlplugin.ActivityInfoMapsRow
	for i := 0; i < numActivities; i++ {
		activity := s.newRandomExecutionActivityRow(shardID, namespaceID, workflowID, runID, rand.Int63())
		activities = append(activities, activity)
	}
	result, err := s.store.ReplaceIntoActivityInfoMaps(activities)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(numActivities, int(rowsAffected))

	filter := &sqlplugin.ActivityInfoMapsFilter{
		ShardID:     shardID,
		NamespaceID: namespaceID,
		WorkflowID:  workflowID,
		RunID:       runID,
	}
	result, err = s.store.DeleteFromActivityInfoMaps(filter)
	s.NoError(err)
	rowsAffected, err = result.RowsAffected()
	s.NoError(err)
	s.Equal(numActivities, int(rowsAffected))

	rows, err := s.store.SelectFromActivityInfoMaps(filter)
	s.NoError(err)
	s.Equal([]sqlplugin.ActivityInfoMapsRow(nil), rows)
}

func (s *historyExecutionActivitySuite) newRandomExecutionActivityRow(
	shardID int32,
	namespaceID primitives.UUID,
	workflowID string,
	runID primitives.UUID,
	scheduleID int64,
) sqlplugin.ActivityInfoMapsRow {
	return sqlplugin.ActivityInfoMapsRow{
		ShardID:      shardID,
		NamespaceID:  namespaceID,
		WorkflowID:   workflowID,
		RunID:        runID,
		ScheduleID:   scheduleID,
		Data:         shuffle.Bytes(testHistoryExecutionActivityData),
		DataEncoding: testHistoryExecutionActivityEncoding,
	}
}
