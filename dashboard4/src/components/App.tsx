"use client"

import React, { useState } from 'react';

interface Row {
    title: string;
    score: number;
}

const rows: Row[] = [
    { title: 'teamA', score: 15 },
    { title: 'teamB', score: 8 },
    { title: 'teamB', score: 9 },
    { title: 'teamC', score: 12 },
];

const App: React.FC = () => {
    const [filters, setFilters] = useState<{ filter: string, field: string }[]>([]);
    const [filteredRows, setFilteredRows] = useState<Row[]>(rows);

    const applyFilter = (data: Row[], filters: { filter: string, field: string }[]): Row[] => {
        const parseExpression = (expression: string, value: string | number): boolean => {
            const operators = /([<>]=?|==|!=)/g;
            const tokens = expression.split(/\s+(and|or)\s+/i);
            let result = false;
            let currentOp = 'or';

            for (const token of tokens) {
                if (token.toLowerCase() === 'and' || token.toLowerCase() === 'or') {
                    currentOp = token.toLowerCase();
                    continue;
                }

                if (typeof value === 'number') {
                    const match = token.match(operators);
                    if (match) {
                        const [operator, number] = [match[0], parseFloat(token.split(match[0])[1])];
                        let comparisonResult = false;

                        switch (operator) {
                            case '>':
                                comparisonResult = value > number;
                                break;
                            case '>=':
                                comparisonResult = value >= number;
                                break;
                            case '<':
                                comparisonResult = value < number;
                                break;
                            case '<=':
                                comparisonResult = value <= number;
                                break;
                            case '==':
                                comparisonResult = value === number;
                                break;
                            case '!=':
                                comparisonResult = value !== number;
                                break;
                        }

                        if (currentOp === 'and') {
                            result = result && comparisonResult;
                        } else {
                            result = result || comparisonResult;
                        }
                    }
                } else {
                    const regex = new RegExp(token.replace(/%/g, '.*'), 'i');
                    const comparisonResult = regex.test(value);

                    if (currentOp === 'and') {
                        result = result && comparisonResult;
                    } else {
                        result = result || comparisonResult;
                    }
                }
            }

            return result;
        };

        return data.filter(row => {
            return filters.every(({ filter, field }) => {
                const value = row[field as keyof Row];
                return parseExpression(filter, value);
            });
        });
    };

    const handleApplyFilter = () => {
        const newFilteredRows = applyFilter(rows, filters);
        setFilteredRows(newFilteredRows);
    };

    const addFilter = () => {
        setFilters([...filters, { filter: '', field: 'title' }]);
    };

    const updateFilter = (index: number, field: string, filter: string) => {
        const newFilters = [...filters];
        newFilters[index] = { field, filter };
        setFilters(newFilters);
    };

    return (
        <div>
            {filters.map((filterObj, index) => (
                <div key={index}>
                    <select value={filterObj.field} onChange={(e) => updateFilter(index, e.target.value, filterObj.filter)}>
                        <option value="title">Title</option>
                        <option value="score">Score</option>
                    </select>
                    <input
                        type="text"
                        value={filterObj.filter}
                        onChange={(e) => updateFilter(index, filterObj.field, e.target.value)}
                        placeholder="Enter filter (e.g., 'team%' or '>10')"
                    />
                </div>
            ))}
            <button onClick={addFilter}>Add Filter</button>
            <button onClick={handleApplyFilter}>Apply Filter</button>
            <ul>
                {filteredRows.map((row, index) => (
                    <li key={index}>{row.title} - {row.score}</li>
                ))}
            </ul>
        </div>
    );
};

export default App;
