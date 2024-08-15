package lang

import (
	"testing"
)

func TestLangConstants(t *testing.T) {
	// Test CmdRootShort and CmdRootLong
	if CmdRootShort != "Watch Auditor" {
		t.Errorf("Expected CmdRootShort to be 'Watch Auditor', but got '%s'", CmdRootShort)
	}

	if CmdRootLong != "Watch Auditor is a tool to audit the watch controller" {
		t.Errorf("Expected CmdRootLong to be 'Watch Auditor is a tool to audit the watch controller', but got '%s'", CmdRootLong)
	}

	if CmdServeShort != "Start the server" {
		t.Errorf("Expected CmdServeShort to be 'Start the server', but got '%s'", CmdServeShort)
	}

	if CmdServeLong != "Start the server" {
		t.Errorf("Expected CmdServeLong to be 'Start the server', but got '%s'", CmdServeLong)
	}

	if SchedulerAuditorSuccessCreation != "Auditor successfully created pod: %s" {
		t.Errorf("Expected SchedulerAuditorSuccessCreation to be 'Auditor successfully created pod', but got '%s'", SchedulerAuditorSuccessCreation)
	}

	if SchedulerAuditorFailedCreation != "Auditor failed to create pod: %s. %s" {
		t.Errorf("Expected SchedulerAuditorFailedCreation to be 'Auditor failed to create pod', but got '%s'", SchedulerAuditorFailedCreation)
	}

	if SchedulerAuditorFailedDeletion != "Auditor failed to delete pod" {
		t.Errorf("Expected SchedulerAuditorFailedDeletion to be 'Auditor failed to delete pod', but got '%s'", SchedulerAuditorFailedDeletion)
	}

	if SchedulerWatcherSuccessDeletion != "Watch Controller successfully deleted pod: %s" {
		t.Errorf("Expected SchedulerWatcherSuccessDeletion to be 'Watch Controller successfully deleted pod, but got '%s'", SchedulerWatcherSuccessDeletion)
	}

	if SchedulerWatcherFailedDeletion != "Watch Controller failed to deleted pod: %s." {
		t.Errorf("Expected SchedulerWatcherFailedDeletion to be 'Watch Controller failed to deleted pod, but got '%s'", SchedulerWatcherFailedDeletion)
	}

	if PromWatchDeletionsName != "watch_controller_deletions_total" {
		t.Errorf("Expected PromWatchDeletionsName to be 'watch_controller_deletions_total', but got '%s'", PromWatchDeletionsName)
	}

	if PromWatchDeletionsHelp != "The total number of watch controller deletions" {
		t.Errorf("Expected PromWatchDeletionsHelp to be 'The total number of watch controller deletions', but got '%s'", PromWatchDeletionsHelp)
	}

	if PromWatchFailuresName != "watch_controller_failures_total" {
		t.Errorf("Expected PromWatchFailuresName to be 'watch_controller_failures_total', but got '%s'", PromWatchFailuresName)
	}

	if PromWatchFailuresHelp != "The total number of watch controller failures" {
		t.Errorf("Expected PromWatchFailuresHelp to be 'The total number of watch controller failures', but got '%s'", PromWatchFailuresHelp)
	}
}
