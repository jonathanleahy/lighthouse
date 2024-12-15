import React, { ComponentType } from 'react';

interface IconWithSlashProps {
    Icon: ComponentType<{ className?: string }>;
    disabled?: boolean;
    className?: string;
}


const IconWithSlash: React.FC<IconWithSlashProps> = ({ Icon, disabled = false, className = "" }) => {
    return (
        <div className="p-2 relative inline-flex items-center justify-center">
            <Icon className={`h-4 w-4 ${disabled ? 'text-muted-foreground' : 'text-foreground'} ${className}`} />
            {disabled && (
                <div className="absolute inset-0 flex items-center justify-center">
                    <div className="w-4 h-px bg-muted-foreground rotate-45 transform origin-center" />
                </div>
            )}
        </div>
    );
};

export default IconWithSlash;