import http from "k6/http";
import { check } from "k6";

export const options = {
  vus: 10,
  duration: "1m",
};
export default function () {
  const res = http.get("http://localhost:8080/cpu");
  check(res, {
    "is status 503": (r) => r.status === 503,
  });
}
