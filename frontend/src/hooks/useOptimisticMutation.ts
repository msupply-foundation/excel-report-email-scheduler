import { MutationConfig, MutationFunction, MutationResultPair, queryCache, QueryKey, useMutation } from 'react-query';

declare type CreateOptimisticValue<TVariables, TOptimisticValue> = (variables: TVariables) => TOptimisticValue;
declare type CreateOptimisticState<TResult, TOptimisticValue> = (
  prevData: TResult | undefined,
  optimisticValue: TOptimisticValue
) => TResult | undefined;

export const useOptimisticMutation = <TResult, TOptimisticValue, TVariables, TSnapshot>(
  queryKey: QueryKey,
  mutationFn: MutationFunction<TResult, TVariables>,
  createOptimisticContent: CreateOptimisticValue<TVariables, TOptimisticValue>,
  createOptimisticState: CreateOptimisticState<TResult, TOptimisticValue>,
  defaultSnapshot: TSnapshot | Promise<TSnapshot>,
  config: MutationConfig<TResult, Error, TVariables, TSnapshot> = {}
): MutationResultPair<TResult, Error, TVariables, TSnapshot> => {
  const refetch = () => queryCache.refetchQueries(queryKey);

  const [callback, status] = useMutation<TResult, Error, TVariables, TSnapshot>(mutationFn, {
    onSuccess: refetch,
    onMutate: (variables: TVariables) => {
      const optimistic = createOptimisticContent(variables);
      queryCache.cancelQueries(queryKey);
      const snapshot = queryCache.getQueryData<TSnapshot | undefined>(queryKey);
      queryCache.setQueryData<TResult | undefined>(queryKey, prevData => {
        return createOptimisticState(prevData, optimistic);
      });
      return snapshot ?? defaultSnapshot;
    },
    onError: (_, __, prevData) => {
      queryCache.setQueryData(queryKey, prevData);
    },
    ...config,
  });

  return [callback, status];
};
