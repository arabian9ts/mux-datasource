# MUX datasource

Visualize data from MUX.  
About MUX: https://www.mux.com/


## Feature
- Query timeseries metrics with Filters
- API Access Token Setting with secure json
- Suggest filters from breakdown result

![query](https://github.com/arabian9ts/mux-datasource/blob/main/src/img/query.png?raw=true)

![config](https://github.com/arabian9ts/mux-datasource/blob/main/src/img/config.png?raw=true)

# MUX API Access Token
You can get the api access token as follows.  
https://docs.mux.com/core/stream-video-files#1-get-an-api-access-token

## Configure the datasource in Grafana
[Add a data source](https://grafana.com/docs/grafana/latest/datasources/add-a-data-source/)

Required fields
1. Access Token ID
2. Access Token Secret

## Configure the data source via provisioning
You can configure this data source with Grafana provisioning system.  
For mode infomation, see [Provisioning Grafana](https://grafana.com/docs/grafana/latest/administration/provisioning/#datasources).

```yml
apiVersion: 1
datasources:
  - name: MUX Data
    type: arabian9ts-mux-datasource
    access: proxy
    orgId: 1
    version: 1
    editable: false
    secureJsonData:
      accessTokenId: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
      accessTokenSecret: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```
