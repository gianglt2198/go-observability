import http from 'k6/http';
import { check, sleep } from 'k6';

// export const options = {
//     stages: [
//         { duration: '1m', target: 1 },
//         { duration: '5m', target: 3 },
//         { duration: '1m', target: 0 },
//     ],

// };

// export default function () {
//     http.get("http://localhost:8081/api/coffee");
//     sleep(1);
// }

export const options = {
    discardResponseBodies: true,
    scenarios: {
        coffees: {
            executor: "per-vu-iterations",
            exec: "coffees",
            vus: 5,
            iterations: 500,
            maxDuration: "1m",
        },
        detail_coffee: {
            executor: "per-vu-iterations",
            exec: "coffee_detail",
            vus: 2,
            iterations: 100,
            maxDuration: '30s',
        }
    }
}

export function coffees() {
    http.get("http://localhost:8081/api/coffee");
}

export function coffee_detail() {
    http.get("http://localhost:8081/api/coffee/1");
}