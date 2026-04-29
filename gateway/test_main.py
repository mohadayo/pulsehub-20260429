from fastapi.testclient import TestClient

from main import app, events_store

client = TestClient(app)


def setup_function():
    events_store.clear()


def test_health():
    r = client.get("/health")
    assert r.status_code == 200
    data = r.json()
    assert data["status"] == "ok"
    assert data["service"] == "gateway"


def test_create_event():
    r = client.post("/events", json={"name": "signup", "payload": {"user": "alice"}})
    assert r.status_code == 200
    data = r.json()
    assert data["id"] == 1
    assert data["name"] == "signup"


def test_create_event_empty_name():
    r = client.post("/events", json={"name": "  "})
    assert r.status_code == 400


def test_list_events():
    client.post("/events", json={"name": "e1"})
    client.post("/events", json={"name": "e2"})
    r = client.get("/events")
    assert r.status_code == 200
    assert r.json()["total"] == 2


def test_get_event():
    client.post("/events", json={"name": "click"})
    r = client.get("/events/1")
    assert r.status_code == 200
    assert r.json()["name"] == "click"


def test_get_event_not_found():
    r = client.get("/events/999")
    assert r.status_code == 404
