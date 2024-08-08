import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '30s', target: 30 },
        { duration: '45s', target: 30 },
        { duration: '30s', target: 50 },
        { duration: '45s', target: 50 },
        { duration: '30s', target: 70 },
        { duration: '45s', target: 70 },
        { duration: '30s', target: 90 },
        { duration: '45s', target: 90 },
        { duration: '30s', target: 110 },
        { duration: '45s', target: 110 },
    ],
};

export default function () {
    const res = http.get('http://app-b.apps.svc.cluster.local:8080/scenarioA');
    check(res, { 'status was 200': (r) => r.status == 200 });
    sleep(5);
}

