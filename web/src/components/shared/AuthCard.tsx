import React from 'react';
import { cn } from "@/lib/utils";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import type { ReactNode } from 'react';

interface AuthCardProps {
  title: string;
  description: string;
  children: ReactNode;
  className?: string;
}

export function AuthCard({ title, description, children, className }: AuthCardProps) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-background p-4">
      <div className={cn("flex flex-col gap-6", "w-full max-w-md min-w-96", className)}>
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl">{title}</CardTitle>
            <CardDescription>{description}</CardDescription>
          </CardHeader>
          <CardContent>
            {children}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}