import ws from 'k6/ws';
import { check } from 'k6';

const duration = __ENV.K6_DURATION || '30s';
const vus = Number(__ENV.K6_VUS || 1);

export const options = {
  scenarios: {
    websocket: {
      executor: 'constant-vus',
      vus,
      duration,
    },
  },
};

export default function () {
  const url = __ENV.REALTIME_GATEWAY_WS_URL;
  if (!url) {
    throw new Error('REALTIME_GATEWAY_WS_URL is required');
  }

  const headers = {};
  if (__ENV.REALTIME_CHANNEL_TICKET) {
    headers.Authorization = `Bearer ${__ENV.REALTIME_CHANNEL_TICKET}`;
  }

  const response = ws.connect(url, { headers }, (socket) => {
    socket.on('open', () => {
      socket.send(JSON.stringify({ type: 'ping' }));
    });
    socket.setTimeout(() => socket.close(), 1000 * 60 * 60);
  });

  check(response, {
    'websocket upgrade succeeded': (result) => result && result.status === 101,
  });
}
