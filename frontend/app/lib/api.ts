import { CreateMonitorRequest, Monitor, MonitorWithLatestResult } from "../types";

export async function fetchMonitorsWithResults(): Promise<MonitorWithLatestResult[]> {
  const res = await fetch("/api/monitors?include=latest_result");
  if (!res.ok) throw new Error(`Failed to fetch monitors: ${res.status}`);
  return res.json();
}

export async function createMonitor(req: CreateMonitorRequest): Promise<Monitor> {
  const res = await fetch("/api/monitors", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `Failed to create monitor: ${res.status}`);
  }
  return res.json();
}

export async function deleteMonitor(id: string): Promise<void> {
  const res = await fetch(`/api/monitors/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) throw new Error(`Failed to delete monitor: ${res.status}`);
}

export async function fetchHealthStatus(): Promise<boolean> {
  try {
    const res = await fetch("/api/health");
    return res.ok;
  } catch {
    return false;
  }
}
