import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export type Metric = string;
export type GroupBy = string;
export type Measurement = string;
export interface MetricFilter {
  group: string;
  operator: string;
  value: string;
};

export interface Query extends DataQuery {
  metricId: Metric;
  measurement: Measurement;
  filters: MetricFilter[];
}

export const DEFAULT_QUERY: Partial<Query> = {
  metricId: 'unique_viewers',
  measurement: 'count',
  filters: [],
};

export interface DataPoint {
  Time: number;
  Value: number;
}

export interface DataSourceResponse {
  datapoints: DataPoint[];
}

export interface DataSourceOptions extends DataSourceJsonData {
}

export interface SecureJsonData {
  tokenId?: string;
  tokenSecret?: string;
}
