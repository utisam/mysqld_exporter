// Scrape `performance_schema.thread_by_user`.

package collector

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

const perfThreadByUserQuery = `
	SELECT processlist_user, COUNT(*)
		FROM performance_schema.threads
		WHERE processlist_user IS NOT NULL GROUP BY (processlist_user);
	`

// Metric descriptors.
var (
	performanceSchemaThreadByUserDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, performanceSchema, "thread_by_user"),
		"The number of thread by user.",
		[]string{"user"}, nil,
	)
)

// ScrapePerfThreadByUser collects from `performance_schema.thread_by_user`.
type ScrapePerfThreadByUser struct{}

// Name of the Scraper. Should be unique.
func (ScrapePerfThreadByUser) Name() string {
	return "perf_schema.thread_by_user"
}

// Help describes the role of the Scraper.
func (ScrapePerfThreadByUser) Help() string {
	return "Collect metrics from performance_schema.threads by user"
}

// Scrape collects data from database connection and sends it over channel as prometheus metric.
func (ScrapePerfThreadByUser) Scrape(db *sql.DB, ch chan<- prometheus.Metric) error {
	// Timers here are returned in picoseconds.
	perfSchemaThreadByUserRows, err := db.Query(perfThreadByUserQuery)
	if err != nil {
		return err
	}
	defer perfSchemaThreadByUserRows.Close()

	var (
		user  string
		count uint64
	)
	for perfSchemaThreadByUserRows.Next() {
		if err := perfSchemaThreadByUserRows.Scan(
			&user, &count,
		); err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(
			performanceSchemaThreadByUserDesc, prometheus.CounterValue, float64(count),
			user,
		)
	}
	return nil
}
