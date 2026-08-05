import { randomUUID } from "node:crypto";
import { Kafka, logLevel, type Consumer, type Producer } from "kafkajs";

export const YjsWorkerCoreCommandTopic = "notezy.yjsworker.core.command.v1";
export const CoreYjsWorkerReplyTopic = "notezy.core.yjsworker.reply.v1";
export const CoreYjsWorkerMaintenanceCommandTopic =
  "notezy.core.yjsworker.maintenance-command.v1";
export const YjsWorkerCoreMaintenanceResultTopic =
  "notezy.yjsworker.core.maintenance-result.v1";

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
  trace: {
    traceParent?: string;
    traceState?: string;
  };
  data: unknown;
};

type CommandEnvelope<D> = {
  schemaVersion: string;
  commandId: string;
  commandType: string;
  blockPackId: string;
  correlationId: string;
  trace: {
    traceParent?: string;
    traceState?: string;
  };
  producer: string;
  occurredAt: string;
  data: D;
};

type ReplyEnvelope<D> = {
  schemaVersion: string;
  commandId: string;
  commandType: string;
  blockPackId: string;
  correlationId: string;
  producer: string;
  respondedAt: string;
  data: D;
  error?: {
    code: string;
    message: string;
    retryable: boolean;
  };
};

type PendingReply = {
  resolve: (data: unknown) => void;
  reject: (error: Error) => void;
  timeout: NodeJS.Timeout;
};

type CoreCommand = {
  commandId: string;
  event: EventEnvelope;
};

function newKafka(): Kafka {
  const brokers = (process.env.KAFKA_BROKERS ?? "127.0.0.1:9094")
    .split(",")
    .map(broker => broker.trim())
    .filter(Boolean);

  return new Kafka({
    clientId: process.env.KAFKA_CLIENT_ID ?? "notezy-yjs-worker",
    brokers,
    logLevel: logLevel.NOTHING,
  });
}

export class CoreCommandProducer {
  private readonly producer: Producer;
  private started = false;

  constructor(kafka: Kafka) {
    this.producer = kafka.producer();
  }

  async start(): Promise<void> {
    if (this.started) return;

    await this.producer.connect();
    this.started = true;
  }

  create<D>(commandType: string, blockPackId: string, data: D): CoreCommand {
    const commandId = randomUUID();
    const correlationId = randomUUID();
    const occurredAt = new Date().toISOString();
    const command: CommandEnvelope<D> = {
      schemaVersion: "v1",
      commandId,
      commandType,
      blockPackId,
      correlationId,
      trace: {},
      producer: "yjs-worker",
      occurredAt,
      data,
    };

    return {
      commandId,
      event: {
        schemaVersion: "v1",
        eventId: commandId,
        eventType: "YjsWorkerCommand",
        aggregateType: "BlockPack",
        aggregateId: blockPackId,
        kafkaKey: blockPackId,
        occurredAt,
        correlationId,
        causationId: commandId,
        trace: {},
        data: command,
      },
    };
  }

  async produce(command: CoreCommand): Promise<void> {
    await this.start();
    await this.producer.send({
      topic: YjsWorkerCoreCommandTopic,
      messages: [{ key: command.event.kafkaKey, value: JSON.stringify(command.event) }],
    });
  }

  async shutdown(): Promise<void> {
    if (!this.started) return;

    await this.producer.disconnect();
    this.started = false;
  }
}

export class CoreReplyConsumer {
  private readonly consumer: Consumer;
  private readonly pendingReplies = new Map<string, PendingReply>();
  private started = false;

  constructor(kafka: Kafka) {
    const instanceId = process.env.HOSTNAME ?? randomUUID();
    this.consumer = kafka.consumer({
      groupId: `notezy-yjsworker-replies-${instanceId}`,
    });
  }

  async start(): Promise<void> {
    if (this.started) return;

    await this.consumer.connect();
    await this.consumer.subscribe({
      topic: CoreYjsWorkerReplyTopic,
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
        const reply = event.data as ReplyEnvelope<unknown>;
        if (
          event.schemaVersion !== "v1" ||
          event.eventType !== "YjsWorkerCommandCompleted" ||
          reply.schemaVersion !== "v1" ||
          typeof reply.commandId !== "string"
        ) {
          return;
        }

        const pendingReply = this.pendingReplies.get(reply.commandId);
        if (pendingReply === undefined) return;

        this.pendingReplies.delete(reply.commandId);
        clearTimeout(pendingReply.timeout);
        if (reply.error !== undefined) {
          pendingReply.reject(
            new Error(`${reply.error.code}: ${reply.error.message}`)
          );

          return;
        }
        pendingReply.resolve(reply.data);
      },
    });
    this.started = true;
  }

  waitForReply<R>(commandId: string, commandType: string): Promise<R> {
    return new Promise<R>((resolve, reject) => {
      const timeout = setTimeout(() => {
        this.pendingReplies.delete(commandId);
        reject(new Error(`YjsWorker command ${commandType} timed out`));
      }, Number(process.env.YJS_WORKER_COMMAND_TIMEOUT_MILLISECONDS ?? 10_000));
      this.pendingReplies.set(commandId, {
        resolve: data => resolve(data as R),
        reject,
        timeout,
      });
    });
  }

  reject(commandId: string, error: unknown): void {
    const pendingReply = this.pendingReplies.get(commandId);
    if (pendingReply === undefined) return;

    this.pendingReplies.delete(commandId);
    clearTimeout(pendingReply.timeout);
    pendingReply.reject(error instanceof Error ? error : new Error(String(error)));
  }

  async shutdown(): Promise<void> {
    for (const pendingReply of this.pendingReplies.values()) {
      clearTimeout(pendingReply.timeout);
      pendingReply.reject(new Error("YjsWorker is shutting down"));
    }
    this.pendingReplies.clear();
    if (!this.started) return;

    await this.consumer.disconnect();
    this.started = false;
  }
}

export class CoreCommandDispatcher {
  private readonly producer: CoreCommandProducer;
  private readonly replyConsumer: CoreReplyConsumer;

  constructor() {
    const kafka = newKafka();
    this.producer = new CoreCommandProducer(kafka);
    this.replyConsumer = new CoreReplyConsumer(kafka);
  }

  async dispatch<D, R>(
    commandType: string,
    blockPackId: string,
    data: D
  ): Promise<R> {
    await this.producer.start();
    await this.replyConsumer.start();

    const command = this.producer.create(commandType, blockPackId, data);
    const reply = this.replyConsumer.waitForReply<R>(command.commandId, commandType);
    try {
      await this.producer.produce(command);
    } catch (error) {
      this.replyConsumer.reject(command.commandId, error);
    }

    return reply;
  }

  async dispatchAsync<D, R>(
    commandType: string,
    blockPackId: string,
    data: D
  ): Promise<{ commandId: string; reply: Promise<R> }> {
    await this.producer.start();
    await this.replyConsumer.start();

    const command = this.producer.create(commandType, blockPackId, data);
    const reply = this.replyConsumer.waitForReply<R>(command.commandId, commandType);
    try {
      await this.producer.produce(command);
    } catch (error) {
      this.replyConsumer.reject(command.commandId, error);
      throw error;
    }

    return { commandId: command.commandId, reply };
  }

  async shutdown(): Promise<void> {
    await this.replyConsumer.shutdown();
    await this.producer.shutdown();
  }
}
