import React, { useState, useCallback } from 'react';
import { QueryEditorProps, toOption } from '@grafana/data';
import { Select, Stack } from '@grafana/ui';
import { AccessoryButton, Space } from '@grafana/experimental';
import { DataSource } from '../datasource';
import { DataSourceOptions, MetricFilter, Query } from '../types';

type Props = QueryEditorProps<
  DataSource,
  Query,
  DataSourceOptions> & {
    groups: string[];
  };

export const FilterList = ({ query, datasource, groups, onChange, onRunQuery }: Props) => {
  const [breakdownGroups, setBreakdownGroups] = useState<string[]>([]);
  const listBreakdowns = useCallback((group: string) => {
    if (group) {
      setBreakdownGroups([]);
      datasource.listBreakdowns(query.metricId, group)
        .then(setBreakdownGroups)
        .catch(console.error);
    }
  }, [datasource, query.metricId]);

  const add = (group: string, operator: string, value: string) => {
    query.filters.push({ group, operator, value });
    onChange(query);
  }

  const remove = (index: number) => {
    query.filters.splice(index, 1);
    onChange(query);
    onRunQuery();
  }

  const update = (index: number, group: string, operator: string, value: string) => {
    query.filters[index] = { group, operator, value };
    onChange(query);
    onRunQuery();
  }

  return (
    <>
      <Stack direction="row" wrap="wrap" gap={3}>
        {query.filters.map((item: MetricFilter, index: number) => (
          <Stack>
            <Select
              width="auto"
              value={item.group ? toOption(item.group) : null}
              showAllSelectedWhenOpen={true}
              options={groups.map((v) => toOption(v))}
              onChange={(e) => { update(index, e.value || '', item.operator, item.value) }}
            />
            <Select
              width="auto"
              value={item.operator ? toOption(item.operator) : null}
              showAllSelectedWhenOpen={true}
              options={['=', '!='].map((v) => toOption(v))}
              onChange={(e) => { update(index, item.group, e.value || '', item.value) }}
            />
            <Select
              width="auto"
              value={item.value ? toOption(item.value) : null}
              showAllSelectedWhenOpen={true}
              onOpenMenu={() => listBreakdowns(item.group)}
              allowCustomValue
              options={breakdownGroups.map((v: string) => toOption(v))}
              onChange={(e) => { update(index, item.group, item.operator, e.value || '') }}
            />
            <AccessoryButton icon="times" variant="secondary" onClick={() => remove(index)} type="button" />
          </Stack>
        ))}
      </Stack>

      <Space v={2} layout='block' />

      <AccessoryButton icon="plus" variant="secondary" onClick={() => { add('', '=', '') }} type="button">
        Add filter
      </AccessoryButton>
    </>
  );
};
