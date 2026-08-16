import { type Consumer, Kafka, logLevel, type Producer } from "kafkajs";
import {
  CoreYjsWorkerMaintenanceCommandTopic,
  YjsWorkerCoreMaintenanceResultTopic,
} from "../../../../../contracts/yjs-worker/v1/yjsworker_contract.js";
import type { YjsCompactionService } from "../../../services/yjs_compaction_service.js";
import type { YjsProjectionService } from "../../../services/yjs_projection_service.js";
import {
  createYjsCompactionResult,
  parseYjsCompactionInput,
} from "../../../types/yjs_compaction.js";
import { parseYjsDocumentState } from "../../../types/yjs_document_state.js";
import { CoreCommandDispatcher } from "../dispatchers/core_command_dispatcher.js";

type EventEnvelope = {
  schemaVersion: string;
  eventId: string;
  eventType: string;
  aggregateType: string;
  aggregateId: string;
  kafkaKey: string;
  occurredAt: string;
  correlationId: string;
  causationId?: string;
  trace?: {
    traceParent?: string;
    traceState?: string;
  };
  data: unknown;
};

type MaintenanceCommand = {
  requestId: string;
  blockPackId: string;
  documentId: string;
  operation: "compact" | "project";
  targetSequence: number;
  correlationId: string;
};

type MaintenanceResult = {
  requestId: string;
  blockPackId: string;
  documentId: string;
  operation: "compact" | "project";
  targetSequence: number;
  success: boolean;
  compactedUntilSequence: number;
  projectedUntilSequence: number;
  error?: string;
};

function newKafka(): Kafka {
  const brokers = (process.env.KAFKA_BROKERS ?? "127.0.0.1:9094")
    .split(",")
    .map(broker => broker.trim())
    .filter(Boolean);

  return new Kafka({
    clientId: process.env.KAFKA_CLIENT_ID ?? "notegic-yjs-worker",
    brokers,
    logLevel: logLevel.NOTHING,
  });
}

export class YjsMaintenanceConsumer {
  private readonly consumer: Consumer;
  private readonly producer: Producer;
  private readonly dispatcher: CoreCommandDispatcher;
  private readonly compactionService: YjsCompactionService;
  private readonly projectionService: YjsProjectionService;
  private started = false;

  constructor(
    dispatcher: CoreCommandDispatcher,
    compactionService: YjsCompactionService,
    projectionService: YjsProjectionService
  ) {
    const kafka = newKafka();
    this.consumer = kafka.consumer({
      groupId: "notegic-yjs-worker-maintenance-v1",
    });
    this.producer = kafka.producer();
    this.dispatcher = dispatcher;
    this.compactionService = compactionService;
    this.projectionService = projectionService;
  }

  async start(): Promise<void> {
    if (this.started) return;

    await this.producer.connect();
    await this.consumer.connect();
    await this.consumer.subscribe({
      topic: CoreYjsWorkerMaintenanceCommandTopic,
      fromBeginning: false,
    });
    await this.consumer.run({
      eachMessage: async ({ message }) => {
        if (message.value === null) return;

        let event: EventEnvelope;
        try {
          event = JSON.parse(message.value.toString("utf8")) as EventEnvelope;
        } catch {
          return;
        }

        const command = this.parseCommand(event);
        if (command === null) return;

        let result: MaintenanceResult;
        try {
          result = await this.execute(command);
        } catch (error) {
          result = {
            requestId: command.requestId,
            blockPackId: command.blockPackId,
            documentId: command.documentId,
            operation: command.operation,
            targetSequence: command.targetSequence,
            success: false,
            compactedUntilSequence: 0,
            projectedUntilSequence: -1,
            error: error instanceof Error ? error.message : String(error),
          };
        }

        await this.producer.send({
          topic: YjsWorkerCoreMaintenanceResultTopic,
          messages: [
            {
              key: command.blockPackId,
              value: JSON.stringify({
                schemaVersion: "v1",
                eventId: command.requestId,
                eventType: "YjsMaintenanceCompleted",
                aggregateType: "BlockPack",
                aggregateId: command.blockPackId,
                kafkaKey: command.blockPackId,
                occurredAt: new Date().toISOString(),
                correlationId: command.correlationId,
                causationId: event.eventId,
                trace: event.trace ?? {},
                data: result,
              } satisfies EventEnvelope),
            },
          ],
        });
      },
    });
    this.started = true;
  }

  private parseCommand(event: EventEnvelope): MaintenanceCommand | null {
    if (
      event.schemaVersion !== "v1" ||
      event.eventType !== "YjsMaintenanceCommand" ||
      event.aggregateType !== "BlockPack" ||
      event.aggregateId !== event.kafkaKey
    ) {
      return null;
    }

    const data = event.data as Partial<MaintenanceCommand>;
    if (
      typeof data.requestId !== "string" ||
      typeof data.blockPackId !== "string" ||
      typeof data.documentId !== "string" ||
      data.blockPackId !== event.aggregateId ||
      (data.operation !== "compact" && data.operation !== "project") ||
      typeof data.targetSequence !== "number" ||
      !Number.isSafeInteger(data.targetSequence) ||
      data.targetSequence < 0 ||
      typeof data.correlationId !== "string"
    ) {
      return null;
    }

    return data as MaintenanceCommand;
  }

  private async execute(
    command: MaintenanceCommand
  ): Promise<MaintenanceResult> {
    if (command.operation === "compact") {
      return this.executeCompaction(command);
    }

    return this.executeProjection(command);
  }

  private async executeCompaction(
    command: MaintenanceCommand
  ): Promise<MaintenanceResult> {
    const loaded = await this.dispatcher.dispatch<
      Record<string, never>,
      { found: boolean; payload?: string }
    >("LoadCompactableYjsDocument", command.blockPackId, {});
    if (!loaded.found || loaded.payload === undefined) {
      return {
        requestId: command.requestId,
        blockPackId: command.blockPackId,
        documentId: command.documentId,
        operation: command.operation,
        targetSequence: command.targetSequence,
        success: true,
        compactedUntilSequence: 0,
        projectedUntilSequence: -1,
      };
    }

    const input = parseYjsCompactionInput(
      Buffer.from(loaded.payload, "base64")
    );
    if (input === null) throw new Error("invalid compactable Yjs document");

    const compacted = this.compactionService.compact(input);
    const resultPayload = createYjsCompactionResult(
      input,
      compacted.snapshot,
      compacted.stateVector
    );
    const applied = await this.dispatcher.dispatch<
      { payload: string },
      { applied: boolean }
    >("ApplyCompactedYjsDocument", command.blockPackId, {
      payload: resultPayload.toString("base64"),
    });

    return {
      requestId: command.requestId,
      blockPackId: command.blockPackId,
      documentId: command.documentId,
      operation: command.operation,
      targetSequence: command.targetSequence,
      success: applied.applied,
      compactedUntilSequence: applied.applied
        ? input.cutoffSequence
        : input.baseCompactedUntilSequence,
      projectedUntilSequence: -1,
    };
  }

  private async executeProjection(
    command: MaintenanceCommand
  ): Promise<MaintenanceResult> {
    const loaded = await this.dispatcher.dispatch<
      Record<string, never>,
      { found: boolean; payload?: string }
    >("LoadYjsDocument", command.blockPackId, {});
    if (!loaded.found || loaded.payload === undefined) {
      return {
        requestId: command.requestId,
        blockPackId: command.blockPackId,
        documentId: command.documentId,
        operation: command.operation,
        targetSequence: command.targetSequence,
        success: true,
        compactedUntilSequence: 0,
        projectedUntilSequence: -1,
      };
    }

    const state = parseYjsDocumentState(Buffer.from(loaded.payload, "base64"));
    if (state === null) throw new Error("invalid Yjs document state");

    const projection = this.projectionService.project({
      blockPackId: command.blockPackId,
      state,
    });
    const applied = await this.dispatcher.dispatch<
      { projection: string },
      { applied: boolean; projectedUntilSequence: number }
    >("ApplyBlockProjection", command.blockPackId, {
      projection: Buffer.from(JSON.stringify(projection)).toString("base64"),
    });

    return {
      requestId: command.requestId,
      blockPackId: command.blockPackId,
      documentId: command.documentId,
      operation: command.operation,
      targetSequence: command.targetSequence,
      success: applied.applied,
      compactedUntilSequence: state.compactedUntilSequence,
      projectedUntilSequence: applied.projectedUntilSequence,
    };
  }

  async shutdown(): Promise<void> {
    if (!this.started) return;

    await this.consumer.disconnect();
    await this.producer.disconnect();
    this.started = false;
  }
}
