import http from "k6/http";
import { check } from "k6";

export const options = {
  vus: 20,
  duration: "30s",
};
export default function () {
  const res = http.get("http://localhost:8080/httpbin");
  check(res, {
    "is status 200": (r) => r.status === 200,
    "is status 503": (r) => r.status === 503,
  });
}
