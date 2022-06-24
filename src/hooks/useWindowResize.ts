import { useState, useEffect } from 'react';

type State = {
  width: undefined | number;
  height: undefined | number;
};

export const useWindowSize = (): State => {
  const [windowSize, setWindowSize] = useState<State>({
    width: undefined,
    height: undefined,
  });

  useEffect(() => {
    const handleResize = () => {
      setWindowSize({
        width: window.innerWidth,
        height: window.innerHeight,
      });
    };

    window.addEventListener('resize', handleResize);
    handleResize();

    return () => window.removeEventListener('resize', handleResize);
  }, []);

  return windowSize;
};
