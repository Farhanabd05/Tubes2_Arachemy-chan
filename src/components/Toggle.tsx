import { useState } from 'react';

interface ToggleProps<T extends string> {
  options: { label: string; value: T }[];
  selected: T;
  onToggle: (value: T) => void;
}

const Toggle = <T extends string>({ options, selected, onToggle }: ToggleProps<T>) => {
  return (
    <div className="flex rounded-md bg-gray-100 p-1 space-x-1">
      {options.map((option) => (
        <button
          key={option.value}
          type="button"
          onClick={() => onToggle(option.value)}
          className={`flex-1 rounded-md px-4 py-2 text-sm font-medium transition-colors
            ${selected === option.value 
              ? 'bg-white text-blue-600 shadow-sm' 
              : 'text-gray-600 hover:bg-gray-50'}`}
        >
          {option.label}
        </button>
      ))}
    </div>
  );
};

export default Toggle;
