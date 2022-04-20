import { useCallback, useState } from 'react';

const useToggle: (initial: boolean) => [boolean, () => void] = (initial: boolean) => {
  const [toggle, setToggle] = useState(initial);

  const onToggle = useCallback(() => {
    setToggle((state) => !state);
  }, [setToggle]);

  return [toggle, onToggle];
};

export { useToggle };
