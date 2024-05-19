import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface Query extends DataQuery {
  metricName: string;
  dimension: string;
}

export const DEFAULT_QUERY: Partial<Query> = {};

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
