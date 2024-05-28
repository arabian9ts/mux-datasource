import React, { useState, useEffect, useCallback } from 'react';
import { QueryEditorProps, toOption } from '@grafana/data';
import { Select, Stack } from '@grafana/ui';
import { AccessoryButton } from '@grafana/experimental';
import { DataSource } from '../datasource';
import { Spacer } from './Spacer';
import { DataSourceOptions, MetricFilter, Query } from '../types';

type Props = QueryEditorProps<
  DataSource,
  Query,
  DataSourceOptions> & {
    groups: string[];
  };

export const FilterList = ({ query, datasource, groups, onChange }: Props) => {
  const [breakdownGroups, setBreakdownGroups] = useState<string[]>([]);
  const [filters, setFilters] = useState<MetricFilter[]>(query.filters);

  useEffect(() => {
    onChange({ ...query, filters });
  }, [filters, onChange, query]);

  const listBreakdowns = useCallback((group: string) => {
    if (group) {
      setBreakdownGroups([]);
      datasource.listBreakdowns(query.metricId, group)
        .then(setBreakdownGroups)
        .catch(console.error);
    }
  }, [datasource, query.metricId]);

  const add = (group: string, operator: string, value: string) => {
    setFilters([...filters, { group, operator, value }]);
  }

  const remove = (index: number) => {
    setFilters(filters.filter((_: any, i: number) => i !== index));
  }

  const update = (index: number, group: string, operator: string, value: string) => {
    setFilters(
      filters.map((item: MetricFilter, i: number) =>
        i === index ? { group, operator, value } : item
      )
    );
  }

  return (
    <>
      <Stack direction="row" wrap="wrap" gap={3}>
        {query.filters.map((item: MetricFilter, index: number) => (
          <Stack key={index}>
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

      <Spacer />

      <AccessoryButton icon="plus" variant="secondary" onClick={() => { add('', '=', '') }} type="button">
        Add filter
      </AccessoryButton>
    </>
  );
};
