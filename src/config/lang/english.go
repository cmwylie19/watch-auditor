package lang

var (
	CmdRootShort  = "Watch Auditor"
	CmdRootLong   = `Watch Auditor is a tool to audit the watch controller`
	CmdServeShort = "Start the server"
	CmdServeLong  = `Start the server`
)

var (
	SchedulerAuditorSuccessCreation = "Auditor successfully created pod: %s"
	SchedulerAuditorFailedCreation  = "Auditor failed to create pod: %s. %s"
	SchedulerAuditorFailedDeletion  = "Auditor failed to delete pod"
	SchedulerWatcherSuccessDeletion = "Watch Controller successfully deleted pod: %s"
	SchedulerWatcherFailedDeletion  = "Watch Controller failed to deleted pod: %s"
)

var (
	PromWatchDeletionsName = "watch_controller_deletions_total"
	PromWatchDeletionsHelp = "The total number of watch controller deletions"
	PromWatchFailuresName  = "watch_controller_failures_total"
	PromWatchFailuresHelp  = "The total number of watch controller failures"
)
