import React from 'react';
import { InlineField, SecretInput, Legend } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { DataSourceOptions, SecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<DataSourceOptions, SecureJsonData> { }

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const { jsonData, secureJsonFields, secureJsonData } = options;
  const handleChange = <Key extends keyof SecureJsonData, Value extends SecureJsonData[Key]>(
    key: Key,
    value: Value,
    set: boolean
  ) => {
    onOptionsChange({
      ...options,
      jsonData,
      secureJsonFields: { ...secureJsonFields, [key]: set },
      secureJsonData: { ...secureJsonData, [key]: value },
    });
  };

  return (
    <>
      <Legend>API Token</Legend>
      <InlineField label='Token ID' labelWidth={16} interactive tooltip={'Secure json field (backend only)'}>
        <SecretInput
          required
          id='config-editor-token-id'
          isConfigured={!!secureJsonFields?.tokenId}
          value={secureJsonData?.tokenId || ''}
          placeholder='Enter your Token ID'
          onReset={() => { handleChange('tokenId', '', false) }}
          onChange={(e) => handleChange('tokenId', e.currentTarget.value, false)}
        />
      </InlineField>
      <InlineField label='Token Secret' labelWidth={16} interactive tooltip={'Secure json field (backend only)'}>
        <SecretInput
          required
          id='config-editor-token-secret'
          isConfigured={!!secureJsonFields?.tokenSecret}
          value={secureJsonData?.tokenSecret || ''}
          placeholder='Enter your Token Secret'
          onReset={() => { handleChange('tokenSecret', '', false) }}
          onChange={(e) => handleChange('tokenSecret', e.currentTarget.value, false)}
        />
      </InlineField>
    </>
  );
}
