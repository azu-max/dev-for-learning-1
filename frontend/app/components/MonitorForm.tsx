"use client";

import { useState } from "react";
import { CreateMonitorRequest } from "../types";
import styles from "./MonitorForm.module.css";

type Props = {
  onSubmit: (req: CreateMonitorRequest) => Promise<void>;
};

export default function MonitorForm({ onSubmit }: Props) {
  const [name, setName] = useState("");
  const [url, setUrl] = useState("");
  const [interval, setInterval] = useState(30);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await onSubmit({ name, url, interval_seconds: interval });
      setName("");
      setUrl("");
      setInterval(30);
    } catch (err) {
      setError(err instanceof Error ? err.message : "追加に失敗しました");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <div className={styles.fields}>
        <input
          className={styles.input}
          type="text"
          placeholder="Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
        <input
          className={styles.input}
          type="url"
          placeholder="https://example.com"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          required
        />
        <div className={styles.intervalGroup}>
          <input
            className={`${styles.input} ${styles.intervalInput}`}
            type="number"
            min={10}
            value={interval}
            onChange={(e) => setInterval(Number(e.target.value))}
          />
          <span className={styles.intervalLabel}>sec</span>
        </div>
        <button className={styles.button} type="submit" disabled={loading}>
          {loading ? "..." : "+ 追加"}
        </button>
      </div>
      {error && <p className={styles.error}>{error}</p>}
    </form>
  );
}
