import http from "k6/http";
import { check } from "k6";

export const options = {
  scenarios: {
    ramping_up_scenario: {
      executor: "ramping-vus",
      startVUs: 10,
      stages: [
        { duration: "1m", target: 10 },
        { duration: "30s", target: 60 },
        { duration: "1m", target: 10 },
      ],
    },
  },
};
export default function () {
  const res = http.get("http://localhost:8080/cpu");
  check(res, {
    "is status 503": (r) => r.status === 503,
  });
}
