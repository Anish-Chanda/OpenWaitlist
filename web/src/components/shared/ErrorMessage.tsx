import React from 'react';

interface ErrorMessageProps {
  message: string;
}

export function ErrorMessage({ message }: ErrorMessageProps) {
  if (!message) return null;
  
  return (
    <div className="text-sm text-destructive bg-destructive/10 p-3 rounded-md">
      {message}
    </div>
  );
}