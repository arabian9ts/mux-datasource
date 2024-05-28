import { DataSourceInstanceSettings, CoreApp } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import {
  Query,
  DataSourceOptions,
  DEFAULT_QUERY,
  Metric,
  GroupBy,
  Measurement,
} from './types';

export class DataSource extends DataSourceWithBackend<Query, DataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<DataSourceOptions>) {
    super(instanceSettings);
  }

  getDefaultQuery(_: CoreApp): Partial<Query> {
    return DEFAULT_QUERY;
  }

  listMetrics(): Promise<Metric[]> {
    return this.getResource('metrics');
  }

  listMeasurements(): Promise<Measurement[]> {
    return this.getResource('measurements');
  }

  listGroups(): Promise<GroupBy[]> {
    return this.getResource('groups');
  }

  listBreakdowns(metricId: Metric, groupBy: GroupBy): Promise<GroupBy[]> {
    return this.getResource(`metrics/${metricId}/breakdowns?groupBy=${groupBy}`);
  }
}
