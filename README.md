# Project scope

This repository project consists of a little tool that connects with the 
[Jira's clockwork API](https://herocoders.atlassian.net/wiki/spaces/CLK/pages/2999975967/Use+the+Clockwork+API),
fetching the registered worklogs on a specified time range. The fetched logs are then grouped by some rules (defined
later on this document), that are aligned with my needs for providing a monthly work hours report for my job.

This tool is autonomous and allows no configuration aside for the time range to fetch, the user creating the report and
the Jira project that should be limited.

# Configuration of the run

The configuration for the program to run is made through environment variables that consists in:

- **CLOCKWORK_API_TOKEN**, which is mandatory, and represents the clockwork API token that allows for the user to fetch
  the logs. This token can be created on
  [Jira -> Apps -> Clockwork
  menu](https://addvolt.atlassian.net/plugins/servlet/ac/clockwork-cloud/clockwork-api-tokens).
- **CLOCKWORK_API_START_DATE**, represents the first date of the worklogs to fetch from the API. This date should follow
  the format **YYYY-MM_DD**. If not provided, it will be set to the first day of the current month.
- **CLOCKWORK_API_END_DATE**, represents the last date of the worklogs to fetch from the API. This date should follow
  the format **YYYY-MM_DD**. If not provided, it will be set to the last day of the current month.
- **CLOCKWORK_API_PROJECT**, which is optional, and represent the project key to be filtered.
- **CLOCKWORK_API_AUTHOR**, which is optional, and represent the user email to that filters the worklog author.
- **CLOCKWORK_PARSER_DEBUG**, which is optional, and can be anything. This enables a more verbose logging of what the
  application is processing.

# Rules for the output reports

This application end task is to generate two reports files, on the CSV format, one for a issue grouping, another
for epic grouping. Both files will be created on the `export` folder alongside the directory where is executed.

## Issue grouping report

The first report will be called `YYYYMM_issues_report.csv` and will list the time spent in each issue. The fields of
this report are:

- **Epic**, which is the epic name of the issue.
- **Issue**, which is the issue summary.
- **Key**, which is the key that represents the user.
- **Type**, which is the type of the issue.
- **HoursSpent**, the hours spent working on the issue.
- **Assignee**, the assignee of the issue.

**Note:** the *YYYYMM* parte of the file name will match the year and month of the starting date of generated report.

## Epic grouping report

The second report will be called `YYYYMM_epic_report.csv` and will list the time spent in each epic. However there are
some exception of this base rule that are considered:

- All **bug type** issues count the working hours in an epic called "Bug", for logging the work on bug issues.
- All issues that are **not assigned** to the author of the worklog are interpreted as **Peer Review** work, and so they
  are grouped in an epic called "Review", for logging peer review work.
- All issues that **don't have an epic** associated will be registered in an epic with its summary.

The fields of this report are:

- **Epic**, which is the epic name of the work done.
- **HoursSpent**, the hours spent working on the epic.

**Note:** the *YYYYMM* part of the file name will match the year and month of the starting date of generated report.
