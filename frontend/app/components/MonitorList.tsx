"use client";

import { MonitorWithLatestResult } from "../types";
import MonitorCard from "./MonitorCard";
import styles from "./MonitorList.module.css";

type Props = {
  monitors: MonitorWithLatestResult[];
  onDelete: (id: string) => Promise<void>;
};

export default function MonitorList({ monitors, onDelete }: Props) {
  if (monitors.length === 0) {
    return (
      <div className={styles.empty}>
        <p>Monitor が登録されていません</p>
        <p className={styles.emptyHint}>上のフォームから追加してください</p>
      </div>
    );
  }

  return (
    <div className={styles.list}>
      {monitors.map((m) => (
        <MonitorCard key={m.id} monitor={m} onDelete={onDelete} />
      ))}
    </div>
  );
}
