# workflow-job-handler

A simple service to handle GitHub [`workflow_job` webhook](https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#workflow_job) and produce a metric.

Designed to be used with [nomad-autoscaler](https://github.com/hashicorp/nomad-autoscaler).

Inspired by [google-workflow-job-to-pubsub](https://github.com/google-github-actions/github-workflow-job-to-pubsub)
