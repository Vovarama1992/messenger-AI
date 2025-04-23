import { Injectable, OnModuleInit } from '@nestjs/common';
import { Kafka, Consumer, Producer } from 'kafkajs';
import { ConfigService } from '@nestjs/config';
import { KafkaTopics } from 'libs/types/src/kafka/kafka-topics.types';

@Injectable()
export class KafkaService implements OnModuleInit {
  private kafka: Kafka;
  private producer: Producer;

  constructor(private readonly configService: ConfigService) {
    this.kafka = new Kafka({
      brokers: [this.configService.get<string>('KAFKA_BROKER')],
    });
  }

  async onModuleInit() {
    this.producer = this.kafka.producer();
    await this.producer.connect();
  }

  async sendMessage<K extends keyof KafkaTopics>(
    topic: K,
    message: KafkaTopics[K],
  ) {
    await this.producer.send({
      topic,
      messages: [{ value: JSON.stringify(message) }],
    });
  }

  getConsumer(): Consumer {
    return this.kafka.consumer({
      groupId: `kafka-consumer-${Math.random()}`,
    });
  }
}
