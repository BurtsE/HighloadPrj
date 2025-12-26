from locust import FastHttpUser, task, between
from datetime import datetime
from locust import FastHttpUser, task, between
from datetime import datetime, UTC
import random

class MetricsUser(FastHttpUser):
    wait_time = between(0.1, 1)

    @task(1)
    def post_metrics(self):
        now = datetime.now(UTC)
        timestamp = now.strftime("%Y-%m-%dT%H:%M:%SZ")
         # Random CPU between 10.0 and 99.9
        cpu = round(random.uniform(10.0, 99.9), 1)
        rps = random.randint(450, 550)
        self.client.post("/metrics", json={
            "timestamp": timestamp,
            "cpu":cpu,
            "rps": rps
        })

    @task(100)
    def get_avg(self):
        self.client.get("/analyze")