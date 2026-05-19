import React from 'react';
import { cn } from '../../lib/utils';

interface CheckboxProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
}

export const Checkbox: React.FC<CheckboxProps> = ({ className, label, id, ...props }) => {
  const inputId = id || `checkbox-${Math.random().toString(36).substr(2, 9)}`;

  return (
    <label htmlFor={inputId} className="flex items-center cursor-pointer select-none">
      <div className="relative">
        <input
          type="checkbox"
          id={inputId}
          className={cn(
            'peer sr-only',
            className
          )}
          {...props}
        />
        <div className="w-5 h-5 border-2 border-gray-300 rounded peer-checked:bg-blue-600 peer-checked:border-blue-600 transition-colors">
          <svg className="w-full h-full text-white opacity-0 peer-checked:opacity-100 p-0.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3">
            <polyline points="20 6 9 17 4 12" />
          </svg>
        </div>
      </div>
      {label && <span className="ml-2 text-sm text-gray-700">{label}</span>}
    </label>
  );
};
