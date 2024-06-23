package lang

var (
	CmdRootShort  = "Watch Auditor"
	CmdRootLong   = `Watch Auditor is a tool to audit the watch controller`
	CmdServeShort = "Start the server"
	CmdServeLong  = `Start the server`
)

var (
	SchedulerWatcherSuccessDeletion    = "Watch Controller successfully deleted pod."
	SchedulerWatcherFailedDeletion     = "Watch Controller failed to deleted pod."
	SchedulerFailedCreate              = "Failed to create neuvector pod: "
	SchedulerFailedDelete              = "Failed to delete neuvector pod: "
	SchedulerWatcherPodSuccessDeletion = "Successfully deleted watcher pod in namespace: "
	SchedulerWatcherPodFailedDeletion  = "Failed to delete watcher pod in namespace: "
)

var (
	PromWatchDeletionsName = "watch_controller_deletions_total"
	PromWatchDeletionsHelp = "The total number of watch controller deletions"
	PromWatchFailuresName  = "watch_controller_failures_total"
	PromWatchFailuresHelp  = "The total number of watch controller failures"
)
