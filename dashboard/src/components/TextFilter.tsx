import React, { useState, ChangeEvent } from 'react';
import { Input } from '@/components/ui/input';

interface TextFilterProps {
    onFilterChange: (filter: string) => void;
}

const TextFilter: React.FC<TextFilterProps> = ({ onFilterChange }) => {
    const [filterText, setFilterText] = useState('');

    const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        setFilterText(value);
        onFilterChange(`%${value}%`); // Using the % pattern for consistency
    };

    return (
        <div className="w-64">
            <Input
                type="text"
                placeholder="Filter data..."
                value={filterText}
                onChange={handleChange}
                className="w-full"
            />
        </div>
    );
};

export default TextFilter;