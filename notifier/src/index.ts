import express, { Request, Response } from "express";

const app = express();
app.use(express.json());

const PORT = process.env.PORT || 8082;
const LOG_LEVEL = process.env.LOG_LEVEL || "info";

interface Notification {
  id: number;
  channel: string;
  message: string;
  sent_at: string;
}

const notifications: Notification[] = [];

function log(level: string, msg: string): void {
  if (level === "debug" && LOG_LEVEL !== "debug") return;
  const ts = new Date().toISOString();
  console.log(`${ts} [${level.toUpperCase()}] notifier: ${msg}`);
}

app.get("/health", (_req: Request, res: Response) => {
  res.json({
    status: "ok",
    service: "notifier",
    timestamp: new Date().toISOString(),
  });
});

app.post("/notify", (req: Request, res: Response) => {
  const { channel, message } = req.body;

  if (!channel || typeof channel !== "string" || !channel.trim()) {
    res.status(400).json({ error: "channel is required" });
    return;
  }
  if (!message || typeof message !== "string" || !message.trim()) {
    res.status(400).json({ error: "message is required" });
    return;
  }

  const notification: Notification = {
    id: notifications.length + 1,
    channel: channel.trim(),
    message: message.trim(),
    sent_at: new Date().toISOString(),
  };
  notifications.push(notification);
  log("info", `Notification sent: id=${notification.id} channel=${notification.channel}`);
  res.json(notification);
});

app.get("/notifications", (_req: Request, res: Response) => {
  log("info", `Listing ${notifications.length} notifications`);
  res.json({ notifications, total: notifications.length });
});

export { app, notifications };

if (require.main === module) {
  app.listen(PORT, () => {
    log("info", `Notifier service started on port ${PORT}`);
  });
}
