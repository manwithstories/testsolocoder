import React from 'react';
import { BookOpen } from 'lucide-react';

interface EmptyProps {
  title: string;
  description?: string;
  action?: React.ReactNode;
  icon?: React.ReactNode;
}

const Empty: React.FC<EmptyProps> = ({ title, description, action, icon }) => {
  return (
    <div className="flex flex-col items-center justify-center py-16 px-4 text-center">
      <div className="w-20 h-20 bg-gray-100 rounded-full flex items-center justify-center mb-4">
        {icon || <BookOpen className="w-10 h-10 text-gray-400" />}
      </div>
      <h3 className="text-lg font-semibold text-gray-900 mb-2">{title}</h3>
      {description && (
        <p className="text-gray-500 max-w-sm mb-6">{description}</p>
      )}
      {action}
    </div>
  );
};

export default Empty;
