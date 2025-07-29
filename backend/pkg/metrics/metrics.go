package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Teams
	TeamsCreatedCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "instancer_teams_created_total",
			Help: "Total number of teams created",
		},
	)

	// Instances
	InstancesCreatedCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "instancer_instance_created_total",
			Help: "Total number of instances created",
		},
		[]string{"challenge"},
	)

	InstancesDeletedCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "instancer_instance_deleted_total",
			Help: "Total number of instances deleted",
		},
	)

	// Flags
	FlagsSubmittedCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "instancer_flags_submitted_total",
			Help: "Total number of flags submitted",
		},
	)

	FlagsCorrectCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "instancer_flags_correct_total",
			Help: "Total number of correct flags submitted",
		},
	)

	FlagsIncorrectCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "instancer_flags_incorrect_total",
			Help: "Total number of incorrect flags submitted",
		},
	)

	FlagsWrongCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "instancer_flags_wrong_total",
			Help: "Total number of wrong flags submitted",
		},
	)

	// Web server
	RequestsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "instancer_requests_total",
			Help: "Total number of requests received",
		},
		[]string{"method", "path", "status"},
	)

	ExceptionsCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "instancer_exceptions_total",
			Help: "Total number of exceptions raised",
		},
	)
)
