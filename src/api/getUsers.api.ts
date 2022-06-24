import { getBackendSrv } from '@grafana/runtime';
import { User } from 'types';

export const getUsers = (datasourceID: number): Promise<User[]> => {
  return getBackendSrv()
    .post('/api/ds/query', {
      queries: [
        {
          datasourceId: datasourceID,
          rawSql: 'SELECT id, name, first_name, last_name, e_mail FROM "user"',
          format: 'table',
        },
      ],
    })
    .then((result) => {
      const frames = result.results.A.frames[0];

      const {
        schema: { fields },
        data: { values },
      } = frames;

      const columnsToExtract = ['id', 'name', 'first_name', 'last_name', 'e_mail'];

      const indexes = fields.reduce((acc: number[], { name }: any, i: number) => {
        if (columnsToExtract.includes(name)) {
          return [...acc, i];
        }
        return acc;
      }, []);

      let results: any = [];

      for (let i = 0; i < values[0].length; i++) {
        const myObj: any = {};
        for (let j = 0; j < indexes.length; j++) {
          myObj[fields[indexes[j]].name] = values[j][i];
        }
        results.push(myObj);
      }

      return results;
    });
};
