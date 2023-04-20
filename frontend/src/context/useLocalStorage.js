import { useState, useEffect } from "react";

export const useLocalStorage = (key) => {
  const [value, setValue] = useState(null);

  useEffect(() => {
    const stored = localStorage.getItem(key);
    const initial = JSON.parse(stored);
    if (initial !== null) {
      setValue(initial);
    }
  }, []);

  useEffect(() => {
    localStorage.setItem(key, JSON.stringify(value));
  }, [value]);

  return [value, setValue];
};
