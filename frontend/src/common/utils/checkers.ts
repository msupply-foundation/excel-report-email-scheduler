export const panelUsesVariable = (sql: string, variableName: string): boolean => {
  return sql.includes(`\${${variableName}}`);
};
