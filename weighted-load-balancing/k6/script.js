import http from "k6/http";
import { check } from "k6";

export const options = {
  scenarios: {
    open_model: {
      executor: "constant-arrival-rate",
      rate: 10,
      timeUnit: "1s",
      duration: "2m",
      preAllocatedVUs: 20,
    },
  },
};
export default function () {
  const res = http.get("http://localhost:8080/get");
  check(res, {
    "is status 200": (r) => r.status === 200,
  });
}
