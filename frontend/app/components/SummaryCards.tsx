"use client";

import { MonitorWithLatestResult } from "../types";
import styles from "./SummaryCards.module.css";

type Props = {
  monitors: MonitorWithLatestResult[];
};

export default function SummaryCards({ monitors }: Props) {
  const healthy = monitors.filter(
    (m) => m.latest_result?.is_healthy === true
  ).length;
  const unhealthy = monitors.filter(
    (m) => m.latest_result !== null && m.latest_result.is_healthy === false
  ).length;
  const total = monitors.length;

  return (
    <div className={styles.container}>
      <div className={`${styles.card} ${styles.healthy}`}>
        <span className={styles.count}>{healthy}</span>
        <span className={styles.label}>Healthy</span>
      </div>
      <div className={`${styles.card} ${styles.unhealthy}`}>
        <span className={styles.count}>{unhealthy}</span>
        <span className={styles.label}>Unhealthy</span>
      </div>
      <div className={`${styles.card} ${styles.total}`}>
        <span className={styles.count}>{total}</span>
        <span className={styles.label}>Total</span>
      </div>
    </div>
  );
}
