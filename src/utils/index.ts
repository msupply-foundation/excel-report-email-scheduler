import { dateTime, DateTime } from '@grafana/data';

export * from './navigation';
export * from './checkers.utils';

export const formatTimeToDate = (time?: string) => {
  const now: DateTime = dateTime(Date.now());
  const d: DateTime = dateTime(now.format('YYYY-MM-DD') + ' ' + time, 'YYYY-MM-DD HH:mm');
  return d.isValid() ? d : undefined;
};
