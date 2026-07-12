import { parsePort } from "./util.js";

export const config = {
  host: process.env.YJS_WORKER_HOST ?? "0.0.0.0",
  port: parsePort(process.env.YJS_WORKER_PORT),
};
