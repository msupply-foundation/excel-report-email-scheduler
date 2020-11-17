import React, { PureComponent } from 'react';

import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './DataSource';
import { MyDataSourceOptions, MyQuery } from './types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  render() {
    return <div className="gf-form"></div>;
  }
}
