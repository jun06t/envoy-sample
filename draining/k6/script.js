import http from "k6/http";
import { check } from "k6";

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: "constant-arrival-rate",
      rate: 1,
      timeUnit: "1s", // 1000 iterations per second, i.e. 1000 RPS
      duration: "30s",
      preAllocatedVUs: 100, // how large the initial pool of VUs would be
      maxVUs: 200, // if the preAllocatedVUs are not enough, we can initialize more
    },
  },
};
export default function () {
  const res = http.get("http://localhost:8080/sleep");
  check(res, {
    "is status 200": (r) => r.status === 200,
  });
}
