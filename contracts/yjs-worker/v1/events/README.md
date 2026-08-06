# YjsWorker event contracts v1

This package owns the YjsWorker command/reply transport metadata and topic
names. The command and reply payload schemas remain in the parent
`contracts/yjs-worker/v1` package. The generic envelope is imported from
`contracts/types/event.go`; no runtime implementation or Kafka client is imported.

This package also owns the Yjs maintenance operation and the Core-to-YjsWorker
command/result contracts. DurableJob requests and Core's maintenance hints are
separate contracts owned by their respective runtimes.

Consumer groups are runtime deployment configuration and are not part of this
contract package.
