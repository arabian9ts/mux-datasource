import React, { useEffect, useState } from 'react';
import { Select } from '@grafana/ui';
import { QueryEditorProps, toOption } from '@grafana/data';
import { FilterList } from './FilterList';
import { DataSource } from '../datasource';
import { EditorField, EditorFieldGroup, EditorRow, EditorRows } from '@grafana/experimental';
import { DataSourceOptions, Query, Metric, GroupBy, Measurement } from '../types';

type Props = QueryEditorProps<DataSource, Query, DataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery, datasource }: Props) {
  const [metricIds, setMetricIds] = useState<Metric[]>([]);
  const [measurements, setMeasurements] = useState<Measurement[]>([]);
  const [groups, setGroups] = useState<GroupBy[]>([]);

  useEffect(() => {
    datasource
      .listMetrics()
      .then((data) => { setMetricIds(data) })
      .catch(console.error);
  }, [datasource])

  useEffect(() => {
    datasource
      .listMeasurements()
      .then((data) => { setMeasurements(data) })
      .catch(console.error);
  }, [datasource])

  useEffect(() => {
    datasource
      .listGroups()
      .then((data) => { setGroups(data) })
      .catch(console.error);
  }, [datasource])

  const handleChange = <Key extends keyof Query, Value extends Query[Key]>(
    key: Key,
    value: Value,
  ) => {
    onChange({...query, [key]: value});
    onRunQuery();
  };

  return (
    <>
      <EditorRows>
        <EditorRow>
          <EditorFieldGroup>
            <EditorField label="Metric" width={26}>
              <Select
                options={metricIds.map((n: string) => toOption(n))}
                value={query.metricId}
                width={28}
                onChange={(e) => handleChange('metricId', e?.value!)}
                className="inline-element"
                isClearable={false}
              />
            </EditorField>
            <EditorField label="Measurement" width={16}>
              <Select
                options={measurements.map((m: string) => toOption(m))}
                value={query.measurement}
                width={28}
                onChange={(e) => handleChange('measurement', e?.value!)}
                className="inline-element"
                isClearable={false}
              />
            </EditorField>
          </EditorFieldGroup>
        </EditorRow>
        <EditorRow>
          <EditorField label="Filters">
            <FilterList
              groups={groups}
              onChange={onChange}
              onRunQuery={onRunQuery}
              datasource={datasource}
              query={query}
            />
          </EditorField>
        </EditorRow>
      </EditorRows>
    </>
  );
}
