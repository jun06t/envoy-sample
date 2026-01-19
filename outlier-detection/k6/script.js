import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  vus: 10,
  duration: "10m",
};

export default function () {
  const res = http.get("http://localhost:10000/");
  check(res, {
    "is status 200": (r) => r.status === 200,
    "is status 5xx": (r) => r.status >= 500,
  });
  sleep(0.1);
}
