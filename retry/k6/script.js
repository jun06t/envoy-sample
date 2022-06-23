import http from "k6/http";
import { check } from "k6";

export const options = {
  vus: 10,
  duration: "30s",
};
export default function () {
  const res = http.get("http://localhost:8080/hello");
  check(res, {
    "is status 200": (r) => r.status === 200,
    "is status 400": (r) => r.status === 400,
    "is status 500": (r) => r.status === 500,
    "is status 503": (r) => r.status === 503,
  });
}
