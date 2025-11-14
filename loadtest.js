import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '30s', target: 10 },
    { duration: '30s', target: 20 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<300'],
    http_req_failed: ['rate<0.001'],
  },
};

const BASE_URL = 'http://localhost:8080';

export default function () {
  const timestamp = Date.now();
  const randomId = Math.floor(Math.random() * 10000);
  const teamName = `team_${timestamp}_${__VU}_${randomId}`;
  
  const teamData = JSON.stringify({
    team_name: teamName,
    members: [
      { 
        user_id: `user1_${timestamp}_${__VU}_${randomId}`, 
        username: `Alice_${timestamp}`, 
        is_active: true 
      },
      { 
        user_id: `user2_${timestamp}_${__VU}_${randomId}`, 
        username: `Bob_${timestamp}`, 
        is_active: true 
      },
      { 
        user_id: `user3_${timestamp}_${__VU}_${randomId}`, 
        username: `Charlie_${timestamp}`, 
        is_active: true 
      }
    ]
  });

  const params = {
    headers: { 'Content-Type': 'application/json' },
  };

  let res = http.post(`${BASE_URL}/team/add`, teamData, params);
  check(res, { 
    'team created': (r) => r.status === 201,
  });

  res = http.get(`${BASE_URL}/team/get?team_name=${teamName}`);
  check(res, { 
    'team retrieved': (r) => r.status === 200,
  });

  sleep(0.2);
}