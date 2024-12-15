"use client"

import { useState, useCallback, useMemo, useEffect } from "react";
import { useRouter } from 'next/navigation';
import { debounce } from "lodash";
import {LucideIcon, RefreshCw} from 'lucide-react';
import { useCustomFields } from '@/lib/customFieldsContext';
import { useFetchData } from "@/hooks/useFetchHook";
import { Card } from "@/components/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Code, LayoutGrid, LayoutList, Edit } from 'lucide-react';
import { Table, TableBody, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { ThemeSwitcher } from "@/components/theme-switcher";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

interface TextFilterProps {
    onFilterChange: (value: string) => void;
    initialValue?: string;
}

type DisplayMode = 'row' | 'card';

interface CardItem {
    id: string;
    label: string;
    value: string;
    isHeader: boolean;
}

interface CardData {
    id: string;
    title: string;
    description: string;
    items: CardItem[];
    chipColor: string;
    chipText: string;
    icon?: LucideIcon;
    repository_name?: string;
    [key: string]: string | boolean | CardItem[] | LucideIcon | undefined;
}

const DEFAULT_CARD_PROPS = {
    description: '',
    chipColor: 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-100',
    chipText: 'Default'
};

interface CustomField {
    id: string;
    name: string;
    type: string;
    isVisibleInCard: boolean;
    isVisibleInRow: boolean;
    sortOrder: 'none' | 'asc' | 'desc';
    filter: string;
    filterEnabled: boolean;
    displayMode: 'row' | 'card';
}

interface ViewConfig {
    defaultView: 'card' | 'table';
    cardConfig: {
        showDescription: boolean;
        showIcon: boolean;
    };
    tableConfig: {
        compact: boolean;
        showIcon: boolean;
    };
}

interface CustomFieldSet {
    id: string;
    name: string;
    fields: CustomField[];
    displayMode: DisplayMode;
    viewConfig: ViewConfig;
}

const DEFAULT_VIEW_CONFIG: ViewConfig = {
    defaultView: 'card',
    cardConfig: {
        showDescription: true,
        showIcon: true
    },
    tableConfig: {
        compact: false,
        showIcon: true
    }
};

interface GridViewProps {
    data: CardData[];
    onSelect: () => void;
    visibleFields: string[];
    viewMode: 'card' | 'table';
    showDescription?: boolean;
    showIcon?: boolean;
    activeSet?: CustomFieldSet;
    sharedLabels?: string[]; // Make it optional
}

interface TableViewProps {
    data: CardData[];
    visibleFields: string[];
    onSelect: () => void;
    compact?: boolean;
    showIcon?: boolean;
    activeSet?: CustomFieldSet;
    sharedLabels?: string[]; // Make it optional
}

interface DashboardControlsProps {
    customFieldSets: CustomFieldSet[];
    activeSetId: string;
    onSetChange: (id: string) => void;
    onEditFields: () => void;
    viewMode: 'card' | 'table';
    onViewModeToggle: () => void;
    onTextFilterChange: (value: string) => void;
    initialTextFilter?: string;
    onRefresh: () => void;
}

interface Repository {
    repository_name: string;
    [key: string]: string | number | boolean;
}

interface FetchedData {
    repositories: Repository[];
}

const TextFilter: React.FC<TextFilterProps> = ({ onFilterChange, initialValue }) => {
    const [value, setValue] = useState(initialValue || '');

    const debouncedFilter = useMemo(
        () => debounce((searchValue: string) => {
            onFilterChange(searchValue);
            localStorage.setItem('dashboardTextFilter', searchValue);
        }, 300),
        [onFilterChange]
    );

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const newValue = e.target.value;
        setValue(newValue);
        debouncedFilter(newValue);
    };

    useEffect(() => {
        return () => {
            debouncedFilter.cancel();
        };
    }, [debouncedFilter]);

    return (
        <Input
            type="text"
            placeholder="Search..."
            value={value}
            onChange={handleChange}
            className="w-64"
        />
    );
};

const DashboardControls: React.FC<DashboardControlsProps> = ({
                                                                 customFieldSets,
                                                                 activeSetId,
                                                                 onSetChange,
                                                                 onEditFields,
                                                                 viewMode,
                                                                 onViewModeToggle,
                                                                 onTextFilterChange,
                                                                 initialTextFilter,
                                                                 onRefresh
                                                             }) => (
    <div className="mb-4 flex justify-between items-center">
        <div className="flex items-center space-x-2">
            <TextFilter onFilterChange={onTextFilterChange} initialValue={initialTextFilter} />
            <Select onValueChange={onSetChange} value={activeSetId}>
                <SelectTrigger className="w-[180px]">
                    <SelectValue placeholder="Select custom set" />
                </SelectTrigger>
                <SelectContent>
                    {customFieldSets.map((set) => (
                        <SelectItem key={set.id} value={set.id}>
                            {set.name}
                        </SelectItem>
                    ))}
                </SelectContent>
            </Select>
            <Button
                variant="ghost"
                size="icon"
                onClick={onEditFields}
            >
                <Edit className="h-4 w-4" />
                <span className="sr-only">Edit custom fields</span>
            </Button>
        </div>
        <div className="flex items-center space-x-2">
            <Button variant="outline" size="sm" onClick={onViewModeToggle}>
                {viewMode === 'card' ? <LayoutList className="mr-2 h-4 w-4" /> : <LayoutGrid className="mr-2 h-4 w-4" />}
                {viewMode === 'card' ? 'Table View' : 'Card View'}
            </Button>
            <Button variant="outline" size="sm" onClick={onRefresh}>
                <RefreshCw className="mr-2 h-4 w-4" />
                Refresh
            </Button>
            <ThemeSwitcher />
        </div>
    </div>
);

const GridView: React.FC<GridViewProps> = ({
                                               data,
                                               onSelect,
                                               visibleFields,
                                               viewMode,
                                               showDescription = true,
                                               showIcon = true,
                                               activeSet,
                                               sharedLabels = [] // Provide default empty array
                                           }) => (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {data.map((item) => {
            const filteredItems = item.items.filter(itemField => {
                const fieldConfig = activeSet?.fields.find(f => f.name === itemField.label);
                return fieldConfig?.isVisibleInCard;
            });

            const cardProps = {
                ...DEFAULT_CARD_PROPS,
                ...item,
                items: filteredItems,
                sharedLabels,
                onSelect,
                visibleFields,
                viewMode,
                showDescription,
                showIcon
            };

            return <Card key={item.id} {...cardProps} />;
        })}
    </div>
);

const TableView: React.FC<TableViewProps> = ({
                                                 data,
                                                 visibleFields,
                                                 onSelect,
                                                 compact = false,
                                                 showIcon = true,
                                                 activeSet,
                                                 sharedLabels = [] // Provide default empty array
                                             }) => (
    <Table className={compact ? 'table-compact' : ''}>
        <TableHeader>
            <TableRow>
                {showIcon && <TableHead></TableHead>}
                <TableHead>Title</TableHead>
                {visibleFields.map(field => (
                    <TableHead key={field}>{field}</TableHead>
                ))}
            </TableRow>
        </TableHeader>
        <TableBody>
            {data.map((item) => {
                const filteredItems = item.items.filter(itemField => {
                    const fieldConfig = activeSet?.fields.find(f => f.name === itemField.label);
                    return fieldConfig?.isVisibleInRow;
                });

                const cardProps = {
                    ...DEFAULT_CARD_PROPS,
                    ...item,
                    items: filteredItems,
                    sharedLabels,
                    onSelect,
                    visibleFields,
                    viewMode: "row" as const,
                    showIcon
                };

                return <Card key={item.id} {...cardProps} />;
            })}
        </TableBody>
    </Table>
);

const getFieldValue = (card: CardData, fieldName: string): string | number | null => {
    const normalizedFieldName = fieldName.toLowerCase().replace(/_/g, '');

    // Check direct properties first
    const directValue = card[fieldName] ??
        card[fieldName.toLowerCase()] ??
        card[normalizedFieldName] ??
        card[fieldName.replace(/_/g, '')];

    if (directValue !== undefined && directValue !== null) {
        return directValue as string | number;
    }

    // Special handling for repository name
    if (normalizedFieldName === 'repositoryname') {
        return card.repository_name || card.title || null;
    }

    // Check items array
    const item = card.items?.find(item =>
        item.label.toLowerCase().replace(/_/g, '') === normalizedFieldName
    );

    return item?.value || null;
};

const parseFilterConditions = (filter: string): string[][] => {
    return filter.split('||')
        .map(orCondition =>
            orCondition.trim()
                .split('&&')
                .map(condition => condition.trim())
                .filter(condition => condition !== '')
        )
        .filter(conditions => conditions.length > 0);
};

const matchesSingleCondition = (
    value: string,
    condition: string,
    fieldType: 'number' | 'text'
): boolean => {
    const normalizedValue = value.toLowerCase();
    const normalizedCondition = condition.toLowerCase();

    if (fieldType === 'number') {
        const numValue = parseFloat(value);
        const operator = condition.match(/^[<>]=?|=/)?.[0] || '=';
        const targetNumber = parseFloat(condition.replace(operator, ''));

        if (!isNaN(numValue) && !isNaN(targetNumber)) {
            switch (operator) {
                case '>': return numValue > targetNumber;
                case '<': return numValue < targetNumber;
                case '>=': return numValue >= targetNumber;
                case '<=': return numValue <= targetNumber;
                default: return numValue === targetNumber;
            }
        }
    }

    // Handle wildcards
    if (normalizedCondition.includes('%')) {
        const pattern = normalizedCondition.replace(/%/g, '');
        const isStartWildcard = condition.startsWith('%');
        const isEndWildcard = condition.endsWith('%');

        if (isStartWildcard && isEndWildcard) {
            return normalizedValue.includes(pattern);
        } else if (isStartWildcard) {
            return normalizedValue.endsWith(pattern);
        } else if (isEndWildcard) {
            return normalizedValue.startsWith(pattern);
        }
    }

    // Exact match
    return normalizedValue === normalizedCondition;
};

const matchesConditions = (
    value: string | number,
    conditions: string[][],
    fieldType: 'number' | 'text'
): boolean => {
    return conditions.some(andConditions =>
        andConditions.every(condition => matchesSingleCondition(String(value), condition, fieldType))
    );
};

const determineFieldType = (value: string | number): 'number' | 'text' => {
    if (typeof value === 'number') return 'number';
    if (typeof value !== 'string') return 'text';
    return !isNaN(parseFloat(value)) && isFinite(Number(value)) ? 'number' : 'text';
};

const getFilteredAndSortedData = (
    cardData: CardData[],
    activeSet: CustomFieldSet | undefined,
    textFilter: string
): CardData[] => {
    if (!activeSet || !cardData) return cardData;

    // First apply text filter if it exists
    let filteredData = cardData;
    const searchTerm = textFilter.toLowerCase().trim();

    if (searchTerm !== '') {
        filteredData = filteredData.filter(card => {
            const searchableValues = [
                card.title,
                card.description,
                ...(card.items?.map(item => item.value) || []),
                // Include any direct property values
                ...Object.values(card).filter(value =>
                    typeof value === 'string' || typeof value === 'number'
                ).map(value => String(value))
            ].filter(Boolean); // Remove null/undefined values

            return searchableValues.some(value =>
                String(value).toLowerCase().includes(searchTerm)
            );
        });
    }

    // Then apply field-specific filters
    const activeFilters = activeSet.fields.filter(field => field.filterEnabled && field.filter);
    if (activeFilters.length > 0) {
        filteredData = filteredData.filter(card => {
            return activeFilters.every(field => {
                const conditions = parseFilterConditions(field.filter);
                const fieldValue = getFieldValue(card, field.name);

                if (!fieldValue) return false;
                return matchesConditions(fieldValue, conditions, determineFieldType(fieldValue));
            });
        });
    }

    // Finally, apply sorting
    const activeSorts = activeSet.fields.filter(field => field.sortOrder !== 'none');
    if (activeSorts.length > 0) {
        filteredData.sort((a, b) => {
            for (const field of activeSorts) {
                const aValue = getFieldValue(a, field.name);
                const bValue = getFieldValue(b, field.name);

                if (aValue !== bValue) {
                    return field.sortOrder === 'asc' ?
                        (aValue === null ? -1 : bValue === null ? 1 : String(aValue).localeCompare(String(bValue))) :
                        (bValue === null ? -1 : aValue === null ? 1 : String(bValue).localeCompare(String(aValue)));
                }
            }
            return 0;

        });
    }

    return filteredData;
};

const getVisibleFields = (set: CustomFieldSet | undefined, currentViewMode: 'card' | 'table'): string[] => {
    if (!set?.fields) return [];

    return set.fields
        .filter(field => {
            if (field.name.toLowerCase() === 'repository_name' || field.name.toLowerCase() === 'description') {
                return false;
            }

            if (currentViewMode === 'card') {
                return field.isVisibleInCard;
            }
            if (currentViewMode === 'table') {
                return field.isVisibleInRow;
            }
            return false;
        })
        .map(field => field.name);
};

export default function DashboardGrid() {
    const { data: fetchedData, loading, error, refetch } = useFetchData<FetchedData>("http://localhost:8083/list-repos");
    const { state, dispatch } = useCustomFields();
    const [textFilter, setTextFilter] = useState(() => {
        if (typeof window !== 'undefined') {
            return localStorage.getItem('dashboardTextFilter') || '';
        }
        return '';
    });
    const router = useRouter();
    const [viewMode, setViewMode] = useState<'card' | 'table'>(() => {
        const activeSet = state.customFieldSets.find(set => set.id === state.activeSetId);
        return activeSet?.displayMode === 'row' ? 'table' : 'card';
    });
    const [forceRefresh, setForceRefresh] = useState(false);

    useEffect(() => {
        const savedCustomFieldSets = localStorage.getItem('customFieldSets');
        const savedActiveSetId = localStorage.getItem('activeSetId');
        if (savedCustomFieldSets) {
            dispatch({ type: 'SET_CUSTOM_FIELD_SETS', payload: JSON.parse(savedCustomFieldSets) });
        }
        if (savedActiveSetId) {
            dispatch({ type: 'SET_ACTIVE_SET_ID', payload: savedActiveSetId });
        }
    }, [dispatch]);

    useEffect(() => {
        const activeSet = state.customFieldSets.find(set => set.id === state.activeSetId);
        if (activeSet?.displayMode) {
            const newViewMode = activeSet.displayMode === 'row' ? 'table' : 'card';
            setViewMode(newViewMode);
        }
    }, [state.activeSetId, state.customFieldSets]);

    const activeSet = useMemo(() =>
            state.customFieldSets.find(set => set.id === state.activeSetId) ||
            state.customFieldSets[0],
        [state.customFieldSets, state.activeSetId]
    );

    const activeViewConfig = useMemo(() =>
            activeSet?.viewConfig || DEFAULT_VIEW_CONFIG,
        [activeSet]
    );

    const processedCardData = useMemo(() => {
        if (!fetchedData?.repositories || !Array.isArray(fetchedData.repositories)) return [];

        return fetchedData.repositories.map((repo: Repository) => {
            const items = Object.entries(repo)
                .map(([key, value]) => ({
                    id: key.toLowerCase().replace(/_/g, '-'),
                    label: key.charAt(0).toUpperCase() + key.slice(1),
                    value: String(value),
                    isHeader: false
                }));

            return {
                ...DEFAULT_CARD_PROPS,
                id: repo.repository_name,
                title: repo.repository_name,
                items: items,
                icon: Code,
                ...Object.fromEntries(
                    Object.entries(repo).map(([key, value]) => [
                        key.toLowerCase().replace(/_/g, ''),
                        value
                    ])
                ),
                ...repo
            };
        });
    }, [fetchedData]);

    const visibleFields = useMemo(() =>
            getVisibleFields(activeSet, viewMode),
        [activeSet, viewMode]
    );

    const filteredAndSortedCardData = useMemo(() => {
        return getFilteredAndSortedData(processedCardData, activeSet, textFilter);
    }, [processedCardData, activeSet, textFilter]);

    const handleSetChange = useCallback((id: string) => {
        dispatch({ type: 'SET_ACTIVE_SET_ID', payload: id });
        localStorage.setItem('activeSetId', id);
    }, [dispatch]);

    const handleViewModeToggle = useCallback(() => {
        setViewMode(prevMode => prevMode === 'card' ? 'table' : 'card');
    }, []);

    const handleTextFilterChange = useCallback((value: string) => {
        setTextFilter(value);
    }, []);

    const handleRefresh = useCallback(() => {
        setForceRefresh(true);
        refetch();
    }, [refetch]);

    useEffect(() => {
        if (forceRefresh) {
            setForceRefresh(false);
        }
    }, [forceRefresh]);

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: {error.message}</div>;
    }

    return (
        <div>
            <DashboardControls
                customFieldSets={state.customFieldSets}
                activeSetId={state.activeSetId}
                onSetChange={handleSetChange}
                onEditFields={() => router.push('/customize-fields')}
                viewMode={viewMode}
                onViewModeToggle={handleViewModeToggle}
                onTextFilterChange={handleTextFilterChange}
                initialTextFilter={textFilter}
                onRefresh={handleRefresh}
            />
            {viewMode === 'card' ? (
                <GridView
                    data={filteredAndSortedCardData}
                    onSelect={() => {}}
                    visibleFields={visibleFields}
                    viewMode={viewMode}
                    showDescription={activeViewConfig.cardConfig.showDescription}
                    showIcon={activeViewConfig.cardConfig.showIcon}
                    activeSet={activeSet}
                />
            ) : (
                <TableView
                    data={filteredAndSortedCardData}
                    visibleFields={visibleFields}
                    onSelect={() => {}}
                    compact={activeViewConfig.tableConfig.compact}
                    showIcon={activeViewConfig.tableConfig.showIcon}
                    activeSet={activeSet}
                />
            )}
        </div>
    );
}