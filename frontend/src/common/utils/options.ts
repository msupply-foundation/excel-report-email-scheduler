import { SelectableValue } from '@grafana/data';
import intl from 'react-intl-universal';

export const getLookbacks = (): Array<SelectableValue<Number>> => [
  { label: intl.get('1day'), value: 24 * 60 * 1000 * 60 },
  { label: intl.get('2days'), value: 2 * 24 * 60 * 1000 * 60 },
  { label: intl.get('3days'), value: 3 * 24 * 60 * 1000 * 60 },
  { label: intl.get('1week'), value: 7 * 24 * 60 * 1000 * 60 },
  { label: intl.get('2weeks'), value: 14 * 24 * 60 * 1000 * 60 },
  { label: intl.get('4weeks'), value: 28 * 24 * 60 * 1000 * 60 },
  { label: intl.get('3months'), value: 28 * 3 * 24 * 60 * 1000 * 60 },
  { label: intl.get('6months'), value: 28 * 6 * 24 * 60 * 1000 * 60 },
  { label: intl.get('1year'), value: 356 * 24 * 60 * 1000 * 60 },
];
