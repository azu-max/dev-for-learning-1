"use client";

import { MonitorWithLatestResult } from "../types";
import styles from "./MonitorCard.module.css";

type Props = {
  monitor: MonitorWithLatestResult;
  onDelete: (id: string) => Promise<void>;
};

function getTimeAgo(checkedAt: string): string {
  const diff = Date.now() - new Date(checkedAt).getTime();
  const seconds = Math.floor(diff / 1000);
  if (seconds < 60) return `${seconds}秒前`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}分前`;
  const hours = Math.floor(minutes / 60);
  return `${hours}時間前`;
}

function getStatusInfo(monitor: MonitorWithLatestResult) {
  const result = monitor.latest_result;
  if (!result) {
    return { icon: "\u26AB", className: styles.unknown, text: "未チェック" };
  }
  if (result.is_healthy) {
    return {
      icon: "\uD83D\uDFE2",
      className: styles.healthy,
      text: `${result.status_code} OK | ${result.response_time}ms | ${getTimeAgo(result.checked_at)}`,
    };
  }
  const errorText = result.error_message
    ? "接続エラー"
    : `${result.status_code} Error`;
  return {
    icon: "\uD83D\uDD34",
    className: styles.unhealthy,
    text: `${errorText} | ${result.response_time}ms | ${getTimeAgo(result.checked_at)}`,
  };
}

export default function MonitorCard({ monitor, onDelete }: Props) {
  const status = getStatusInfo(monitor);

  return (
    <div className={`${styles.card} ${status.className}`}>
      <div className={styles.header}>
        <div className={styles.info}>
          <span className={styles.icon}>{status.icon}</span>
          <span className={styles.name}>{monitor.name}</span>
          <span className={styles.url}>{monitor.url}</span>
        </div>
        <button
          className={styles.deleteButton}
          onClick={() => onDelete(monitor.id)}
          title="削除"
        >
          削除
        </button>
      </div>
      <div className={styles.status}>{status.text}</div>
    </div>
  );
}
