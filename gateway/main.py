import logging
import os
from datetime import datetime, timezone

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

logging.basicConfig(
    level=os.getenv("LOG_LEVEL", "INFO"),
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
)
logger = logging.getLogger("gateway")

app = FastAPI(title="PulseHub API Gateway")

ANALYTICS_URL = os.getenv("ANALYTICS_URL", "http://localhost:8081")
NOTIFIER_URL = os.getenv("NOTIFIER_URL", "http://localhost:8082")

events_store: list[dict] = []


class Event(BaseModel):
    name: str
    payload: dict | None = None


@app.get("/health")
def health():
    return {"status": "ok", "service": "gateway", "timestamp": datetime.now(timezone.utc).isoformat()}


@app.post("/events")
def create_event(event: Event):
    if not event.name.strip():
        raise HTTPException(status_code=400, detail="Event name must not be empty")
    record = {
        "id": len(events_store) + 1,
        "name": event.name,
        "payload": event.payload,
        "created_at": datetime.now(timezone.utc).isoformat(),
    }
    events_store.append(record)
    logger.info("Event created: id=%d name=%s", record["id"], record["name"])
    return record


@app.get("/events")
def list_events():
    logger.info("Listing %d events", len(events_store))
    return {"events": events_store, "total": len(events_store)}


@app.get("/events/{event_id}")
def get_event(event_id: int):
    for ev in events_store:
        if ev["id"] == event_id:
            return ev
    raise HTTPException(status_code=404, detail=f"Event {event_id} not found")
