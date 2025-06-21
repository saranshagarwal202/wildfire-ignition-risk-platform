package kafka

const (
	TopicDownloadTasks    = "download.tasks"
	TopicProcessTasks     = "process.tasks"
	TopicAnalyticsResults = "analytics.results"
)

var RequiredTopics = []string{
	TopicDownloadTasks,
	TopicProcessTasks,
	TopicAnalyticsResults,
}
