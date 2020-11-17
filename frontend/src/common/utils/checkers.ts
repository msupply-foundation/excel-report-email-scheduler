export const panelUsesVariable = (sql: string, variableName: string): boolean => {
  return sql.includes(`\${${variableName}`);
};

export const panelUsesMacro = (sql: string): boolean => {
  const timeToRegEx = RegExp(/\$__timeTo()/g);
  const timeFromRegEx = RegExp(/\$__timeFrom()/g);
  const timeFilterRegEx = RegExp(/\$__timeFilter\([a-zA-Z]+\)/g);

  return timeToRegEx.test(sql) || timeFromRegEx.test(sql) || timeFilterRegEx.test(sql);
};

export const panelUsesUnsupportedMacro = (sql: string) => {
  // Match all macros starting with $__ except supported ones: timeFrom, timeTo, timeFilter
  const regexp = RegExp(/(?!.*\$__timeFrom\(\).*)(?!.*\$__timeTo\(\).*)(?!.*\$__timeFilter\(.+\).*)(\$__.*)/g);
  return regexp.test(sql);
};
