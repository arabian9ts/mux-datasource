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
      <Legend>API Access Token</Legend>
      <InlineField label='Access Token ID' labelWidth={22} interactive tooltip={'Secure json field (backend only)'}>
        <SecretInput
          required
          id='config-editor-access-token-id'
          width={30}
          isConfigured={!!secureJsonFields?.tokenId}
          value={secureJsonData?.tokenId || ''}
          placeholder='Enter your Access Token ID'
          onReset={() => { handleChange('tokenId', '', false) }}
          onChange={(e) => handleChange('tokenId', e.currentTarget.value, false)}
        />
      </InlineField>
      <InlineField label='Access Token Secret' labelWidth={22} interactive tooltip={'Secure json field (backend only)'}>
        <SecretInput
          required
          id='config-editor-access-token-secret'
          width={30}
          isConfigured={!!secureJsonFields?.tokenSecret}
          value={secureJsonData?.tokenSecret || ''}
          placeholder='Enter your Access Token Secret'
          onReset={() => { handleChange('tokenSecret', '', false) }}
          onChange={(e) => handleChange('tokenSecret', e.currentTarget.value, false)}
        />
      </InlineField>
    </>
  );
}
