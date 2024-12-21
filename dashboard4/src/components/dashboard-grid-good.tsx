"use client"

import { useState, useCallback, useMemo, useEffect } from "react"
import { Card } from "@/components/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Code, LayoutGrid, LayoutList, Edit } from 'lucide-react'
import { Table, TableBody, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { ThemeSwitcher } from "@/components/theme-switcher"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { useRouter } from 'next/navigation'
import { useCustomFields } from '@/lib/customFieldsContext'
import { useFetchData } from "@/hooks/useFetchHook"
import { debounce } from "lodash"
import { LucideIcon } from 'lucide-react'

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
    sharedLabels: string[];
    onSelect: () => void;
    visibleFields: string[];
    viewMode: 'card' | 'table';
    showDescription?: boolean;
    showIcon?: boolean;
    activeSet?: CustomFieldSet;
}

interface TableViewProps {
    data: CardData[];
    visibleFields: string[];
    onSelect: () => void;
    sharedLabels: string[];
    compact?: boolean;
    showIcon?: boolean;
    activeSet?: CustomFieldSet;
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
}

interface Repository {
    repository_name: string;
    [key: string]: string | number | boolean;  // Allow for dynamic fields
}

interface FetchedData {
    repositories: Repository[];
}

// interface FetchedData {
//     repositories: Repository[];
// }

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

const cardData: CardData[] = [
    {
        id: "project-alpha",
        title: "Project Alpha",
        description: "Ongoing development of new features",
        items: [
            { id: "tasks", label: "Tasks", value: "23", isHeader: true },
            { id: "due", label: "Due", value: "5d", isHeader: true },
            { id: "budget", label: "Budget", value: "$15k", isHeader: true },
            { id: "team", label: "Team", value: "6", isHeader: true },
            { id: "status", label: "Status", value: "On Track", isHeader: true },
            { id: "priority", label: "Priority", value: "High", isHeader: true },
            { id: "meeting", label: "Meeting", value: "Mon 10AM", isHeader: false },
            { id: "milestone", label: "Milestone", value: "Feature X", isHeader: false },
            { id: "blocker", label: "Blocker", value: "None", isHeader: false },
            { id: "feedback", label: "Feedback", value: "Pending", isHeader: false },
        ],
        chipColor: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100",
        chipText: "Active",
        icon: Code,
    }
];

const sharedLabels: string[] = ["Meeting", "Milestone", "Blocker", "Feedback"];

const DashboardControls: React.FC<DashboardControlsProps> = ({
                                                                 customFieldSets,
                                                                 activeSetId,
                                                                 onSetChange,
                                                                 onEditFields,
                                                                 viewMode,
                                                                 onViewModeToggle,
                                                                 onTextFilterChange,
                                                                 initialTextFilter
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
            <ThemeSwitcher />
        </div>
    </div>
);

const GridView: React.FC<GridViewProps> = ({
                                               data,
                                               sharedLabels,
                                               onSelect,
                                               visibleFields,
                                               viewMode,
                                               showDescription = true,
                                               showIcon = true,
                                               activeSet
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
                                                 sharedLabels,
                                                 compact = false,
                                                 showIcon = true,
                                                 activeSet
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

const getFilteredAndSortedData = (
    cardData: CardData[],
    activeSet: CustomFieldSet | undefined,
    textFilter: string
): CardData[] => {
    if (!activeSet || !cardData) return cardData;

    const hasActiveFilters = activeSet.fields.some(field => field.filterEnabled && field.filter);
    const hasActiveSorts = activeSet.fields.some(field => field.sortOrder !== 'none');
    const hasTextFilter = textFilter && textFilter.trim() !== '';

    if (!hasActiveFilters && !hasActiveSorts && !hasTextFilter) return cardData;

    return cardData
        .filter(card => {
            if (hasTextFilter) {
                const searchTerm = textFilter.toLowerCase().replace(/%/g, '');
                const matchesText =
                    card.title?.toLowerCase().includes(searchTerm) ||
                    card.description?.toLowerCase().includes(searchTerm) ||
                    card.items?.some(item =>
                        item.value?.toLowerCase().includes(searchTerm)
                    );

                if (!matchesText) return false;
            }

            if (!hasActiveFilters) return true;

            return activeSet.fields.every(field => {
                if (!field.filterEnabled || !field.filter) return true;

                const normalizedFieldName = field.name.toLowerCase().replace(/_/g, '');
                const directValues = [
                    card[field.name],
                    card[field.name.toLowerCase()],
                    card[normalizedFieldName],
                    card[field.name.replace(/_/g, '')],
                    field.name.toLowerCase() === 'repository_name' ? card.repository_name : null,
                    field.name.toLowerCase() === 'repository_name' ? card.title : null
                ].filter(Boolean);

                const sampleValue = directValues.find(v => v !== undefined && v !== null);
                const fieldType = determineFieldType(sampleValue as string | number);
                const conditions = parseFilterConditions(field.filter);

                for (const value of directValues) {
                    if (value && matchesConditions(String(value), conditions, fieldType)) {
                        return true;
                    }
                }

                const matchingItems = card.items?.filter(item => {
                    const itemLabelNormalized = item.label.toLowerCase().replace(/_/g, '');
                    return itemLabelNormalized === normalizedFieldName;
                });

                return matchingItems?.some(item => {
                    const itemType = determineFieldType(item.value);
                    return matchesConditions(String(item.value), conditions, itemType);
                }) || false;
            });
        })
        .sort((a, b) => {
            if (!hasActiveSorts) return 0;

            for (const field of activeSet.fields) {
                if (field.sortOrder !== 'none') {
                    const aValue = getSortValue(a, field.name);
                    const bValue = getSortValue(b, field.name);

                    if (aValue !== bValue) {
                        return field.sortOrder === 'asc'
                            ? compareValues(aValue, bValue)
                            : compareValues(bValue, aValue);
                    }
                }
            }
            return 0;
        });
};

const parseFilterConditions = (filter: string): string[][] => {
    if (!filter) return [];
    return filter.split('||').map(orCondition =>
        orCondition.trim().split('&&').map(condition => condition.trim())
    );
};

const determineFieldType = (value: string | number): 'number' | 'text' => {
    if (typeof value === 'number') return 'number';
    if (typeof value !== 'string') return 'text';
    return !isNaN(parseFloat(value)) && isFinite(Number(value)) ? 'number' : 'text';
};

const matchesConditions = (value: string, conditions: string[][], fieldType: 'number' | 'text'): boolean => {
    return conditions.some(andConditions =>
        andConditions.every(condition => matchesSingleCondition(value, condition, fieldType))
    );
};

const matchesSingleCondition = (value: string, condition: string, fieldType: 'number' | 'text'): boolean => {
    if (fieldType === 'number') {
        const operator = condition.match(/^[<>]=?|=/)?.[0];
        const number = parseFloat(condition.replace(operator || '', ''));

        if (!isNaN(number)) {
            const numValue = parseFloat(value);
            switch (operator) {
                case '>': return numValue > number;
                case '<': return numValue < number;
                case '>=': return numValue >= number;
                case '<=': return numValue <= number;
                case '=': return numValue === number;
                default: return numValue === number;
            }
        }
    }

    if (condition.includes('%')) {
        const isStartWildcard = condition.startsWith('%');
        const isEndWildcard = condition.endsWith('%');
        const cleanPattern = condition.replace(/%/g, '').toLowerCase();
        const stringValue = String(value).toLowerCase();

        if (isStartWildcard && isEndWildcard) {
            return stringValue.includes(cleanPattern);
        } else if (isStartWildcard) {
            return stringValue.endsWith(cleanPattern);
        } else if (isEndWildcard) {
            return stringValue.startsWith(cleanPattern);
        }
    }

    return String(value).toLowerCase() === condition.toLowerCase();
};

const getSortValue = (card: CardData, fieldName: string): string | number => {
    const normalizedFieldName = fieldName.toLowerCase().replace(/_/g, '');

    const directValue = card[fieldName] ||
        card[fieldName.toLowerCase()] ||
        card[normalizedFieldName] ||
        card[fieldName.replace(/_/g, '')];

    if (directValue !== undefined && directValue !== null) {
        return directValue as string | number;
    }

    const item = card.items?.find(item =>
        item.label.toLowerCase().replace(/_/g, '') === normalizedFieldName
    );

    return item?.value || '';
};

const compareValues = (a: string | number, b: string | number): number => {
    if (typeof a === 'number' && typeof b === 'number') {
        return a - b;
    }
    return String(a).localeCompare(String(b));
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
    const { data: fetchedData, loading, error } = useFetchData<FetchedData>("http://localhost:8083/list-repos");
    const { state, dispatch } = useCustomFields()
    const [textFilter, setTextFilter] = useState(() => {
        if (typeof window !== 'undefined') {
            return localStorage.getItem('dashboardTextFilter') || '';
        }
        return '';
    });
    const router = useRouter()

    const [viewMode, setViewMode] = useState<'card' | 'table'>(() => {
        const activeSet = state.customFieldSets.find(set => set.id === state.activeSetId);
        return activeSet?.displayMode === 'row' ? 'table' : 'card';
    });

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
        if (!fetchedData?.repositories || !Array.isArray(fetchedData.repositories)) return cardData;

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

    const filteredAndSortedCardData = useMemo(() =>
            getFilteredAndSortedData(processedCardData, activeSet, textFilter),
        [processedCardData, activeSet, textFilter]
    );

    useEffect(() => {
        if (state.customFieldSets.length === 0 && fetchedData?.repositories) {
            const allFields = [...new Set(
                fetchedData.repositories.flatMap(repo =>
                    Object.keys(repo)
                )
            )]
                .map(key => key.charAt(0).toUpperCase() + key.slice(1));

            const defaultSet: CustomFieldSet = {
                id: '1',
                name: 'Default Fields',
                fields: allFields.map((field, index) => ({
                    id: `${index + 1}`,
                    name: field,
                    type: 'text',
                    isVisibleInCard: true,
                    isVisibleInRow: true,
                    sortOrder: 'none',
                    filter: '',
                    filterEnabled: false,
                    displayMode: 'row'
                })),
                displayMode: 'card',
                viewConfig: DEFAULT_VIEW_CONFIG  // Include viewConfig in the default set
            };
            dispatch({ type: 'SET_CUSTOM_FIELD_SETS', payload: [defaultSet] });
            dispatch({ type: 'SET_ACTIVE_SET_ID', payload: defaultSet.id });
        }
    }, [fetchedData, state.customFieldSets.length, dispatch]);

    const handleCardSelect = useCallback(() => {
        // Functionality to be implemented if needed
        console.log('Card selection functionality to be implemented');
    }, []);

    const handleSetChange = useCallback((setId: string) => {
        dispatch({ type: 'SET_ACTIVE_SET_ID', payload: setId });
        localStorage.setItem('activeSetId', setId);
    }, [dispatch]);

    const toggleViewMode = useCallback(() => {
        const newMode = viewMode === 'card' ? 'table' : 'card';
        setViewMode(newMode);

        const updatedSets = state.customFieldSets.map(set => ({
            ...set,
            displayMode: newMode === 'table' ? ('row' as DisplayMode) : ('card' as DisplayMode),
            viewConfig: {
                ...DEFAULT_VIEW_CONFIG,
                ...(set.viewConfig || {})
            }
        }));

        dispatch({ type: 'SET_CUSTOM_FIELD_SETS', payload: updatedSets });
        localStorage.setItem('customFieldSets', JSON.stringify(updatedSets));
    }, [viewMode, state.customFieldSets, dispatch]);

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error: {error.message}</div>;

    return (
        <div>
            <DashboardControls
                customFieldSets={state.customFieldSets}
                activeSetId={state.activeSetId}
                onSetChange={handleSetChange}
                onEditFields={() => router.push('/customize-fields')}
                viewMode={viewMode}
                onViewModeToggle={toggleViewMode}
                onTextFilterChange={setTextFilter}
                initialTextFilter={textFilter}
            />

            {viewMode === 'card' ? (
                <GridView
                    data={filteredAndSortedCardData}
                    sharedLabels={sharedLabels}
                    onSelect={handleCardSelect}
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
                    onSelect={handleCardSelect}
                    sharedLabels={sharedLabels}
                    compact={activeViewConfig.tableConfig.compact}
                    showIcon={activeViewConfig.tableConfig.showIcon}
                    activeSet={activeSet}
                />
            )}
        </div>
    );
}