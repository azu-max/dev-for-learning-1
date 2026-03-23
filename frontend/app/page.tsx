"use client";

import { useCallback, useEffect, useState } from "react";
import { MonitorWithLatestResult, CreateMonitorRequest } from "./types";
import {
  fetchMonitorsWithResults,
  createMonitor,
  deleteMonitor,
  fetchHealthStatus,
} from "./lib/api";
import SummaryCards from "./components/SummaryCards";
import MonitorForm from "./components/MonitorForm";
import MonitorList from "./components/MonitorList";
import styles from "./page.module.css";

const POLL_INTERVAL = 30_000;

export default function Home() {
  const [monitors, setMonitors] = useState<MonitorWithLatestResult[]>([]);
  const [apiHealthy, setApiHealthy] = useState<boolean | null>(null);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);

  const refresh = useCallback(async () => {
    try {
      const [data, health] = await Promise.all([
        fetchMonitorsWithResults(),
        fetchHealthStatus(),
      ]);
      setMonitors(data);
      setApiHealthy(health);
      setLastUpdated(new Date());
    } catch {
      setApiHealthy(false);
    }
  }, []);

  useEffect(() => {
    refresh();
    const id = window.setInterval(refresh, POLL_INTERVAL);
    return () => clearInterval(id);
  }, [refresh]);

  const handleCreate = async (req: CreateMonitorRequest) => {
    await createMonitor(req);
    await refresh();
  };

  const handleDelete = async (id: string) => {
    await deleteMonitor(id);
    await refresh();
  };

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1 className={styles.title}>Health Check Monitor</h1>
        <span
          className={`${styles.apiStatus} ${apiHealthy === true ? styles.apiOk : apiHealthy === false ? styles.apiDown : ""}`}
        >
          API: {apiHealthy === true ? "ok" : apiHealthy === false ? "down" : "..."}
        </span>
      </header>

      <main className={styles.main}>
        <SummaryCards monitors={monitors} />
        <MonitorForm onSubmit={handleCreate} />
        <MonitorList monitors={monitors} onDelete={handleDelete} />
      </main>

      <footer className={styles.footer}>
        <span>30秒ごとに自動更新</span>
        {lastUpdated && (
          <span>最終更新: {lastUpdated.toLocaleTimeString("ja-JP")}</span>
        )}
      </footer>
    </div>
  );
}
