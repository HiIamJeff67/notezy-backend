# YjsWorker event contracts v1

This package owns the YjsWorker command/reply transport metadata and topic
names. The command and reply payload schemas remain in the parent
`contracts/yjsworker/v1` package. The generic envelope is imported from
`contracts/events`; no runtime implementation or Kafka client is imported.
