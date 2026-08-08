import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: Number(__ENV.K6_VUS || 1),
  duration: __ENV.K6_DURATION || '5m',
};

export default function () {
  const metricsURL = __ENV.KAFKA_METRICS_URL;
  if (!metricsURL) {
    throw new Error('KAFKA_METRICS_URL is required');
  }

  const response = http.get(metricsURL);
  check(response, {
    'metrics endpoint is available': (result) => result.status === 200,
    'consumer lag stays below threshold': (result) => {
      const threshold = Number(__ENV.KAFKA_LAG_THRESHOLD || 1000);
      const lag = Number(result.headers['X-Kafka-Consumer-Lag'] || 0);
      return lag <= threshold;
    },
  });
}
