import http from 'k6/http';
import { check } from 'k6';

export default function() {
    const res = http.post(
        'http://localhost:8080',
        JSON.stringify({ text: 'ttt' }),
        {
            headers: { 'Content-Type': 'application/json' },
        },
    );

    check(res, {
        'is status 200': (r) => r.status === 200,
    });
}
