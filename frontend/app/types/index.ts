export type Monitor = {
  id: string;
  name: string;
  url: string;
  interval_seconds: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

export type CheckResult = {
  id: string;
  monitor_id: string;
  status_code: number;
  response_time: number;
  is_healthy: boolean;
  error_message: string;
  checked_at: string;
};

export type MonitorWithLatestResult = Monitor & {
  latest_result: CheckResult | null;
};

export type CreateMonitorRequest = {
  name: string;
  url: string;
  interval_seconds: number;
};
