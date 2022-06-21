import http from "k6/http";
import { check } from "k6";

export const options = {
  scenarios: {
    ramping_up_scenario: {
      executor: "ramping-vus",
      startVUs: 5,
      stages: [
        { duration: "3m", target: 5 },
        { duration: "1m", target: 40 },
        { duration: "3m", target: 5 },
      ],
    },
  },
};
export default function () {
  const res = http.get("http://localhost:8080/cpu");
  check(res, {
    "is status 200": (r) => r.status === 200,
    "is status 503": (r) => r.status === 503,
  });
}
