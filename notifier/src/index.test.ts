import request from "supertest";
import { app, notifications } from "./index";

beforeEach(() => {
  notifications.length = 0;
});

describe("GET /health", () => {
  it("returns ok status", async () => {
    const res = await request(app).get("/health");
    expect(res.status).toBe(200);
    expect(res.body.status).toBe("ok");
    expect(res.body.service).toBe("notifier");
  });
});

describe("POST /notify", () => {
  it("creates a notification", async () => {
    const res = await request(app)
      .post("/notify")
      .send({ channel: "email", message: "Hello world" });
    expect(res.status).toBe(200);
    expect(res.body.id).toBe(1);
    expect(res.body.channel).toBe("email");
    expect(res.body.message).toBe("Hello world");
  });

  it("rejects empty channel", async () => {
    const res = await request(app)
      .post("/notify")
      .send({ channel: "", message: "test" });
    expect(res.status).toBe(400);
  });

  it("rejects missing message", async () => {
    const res = await request(app)
      .post("/notify")
      .send({ channel: "slack" });
    expect(res.status).toBe(400);
  });
});

describe("GET /notifications", () => {
  it("lists notifications", async () => {
    await request(app).post("/notify").send({ channel: "sms", message: "hi" });
    await request(app).post("/notify").send({ channel: "email", message: "bye" });
    const res = await request(app).get("/notifications");
    expect(res.status).toBe(200);
    expect(res.body.total).toBe(2);
    expect(res.body.notifications).toHaveLength(2);
  });
});
