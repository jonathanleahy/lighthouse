import React from 'react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Box } from 'lucide-react';

interface Item {
    id: string;
    label: string;
    value: string;
    isHeader?: boolean;
    items?: Item[];
}

interface CardComponentProps {
    id: string;
    title: string;
    description: string;
    items: Item[];
    chipColor: string;
    chipText: string;
    visibleFields: string[];
    viewMode: string;
    showDescription?: boolean;
    showIcon?: boolean;
    sharedLabels: string[];
    onSelect: () => void;
}

interface GridViewProps {
    data: Item[];
    sharedLabels: string[];
    onSelect: () => void;
    visibleFields: string[];
    viewMode: string;
    showDescription?: boolean;
    showIcon?: boolean;
    activeSet?: {
        fields: Array<{
            name: string;
            isVisibleInRow: boolean;
        }>;
    };
}

interface TableViewProps {
    data: Item[];
    visibleFields: string[];
    onSelect: () => void;
    sharedLabels: string[];
    compact?: boolean;
    showIcon?: boolean;
    activeSet?: {
        fields: Array<{
            name: string;
            isVisibleInRow: boolean;
        }>;
    };
}

const CardComponent: React.FC<CardComponentProps> = ({
                                                         id,
                                                         title,
                                                         description,
                                                         items,
                                                         chipColor,
                                                         chipText,
                                                         visibleFields,
                                                         viewMode,
                                                         showDescription = true,
                                                         showIcon = true
                                                     }) => {
    const handleClick = (e: React.MouseEvent<HTMLDivElement | HTMLTableRowElement, MouseEvent>) => {
        if (!(e.target as HTMLElement).closest('.checkbox-container')) {
            window.location.href = `/microservice/${id}`;
        }
    };

    if (viewMode === 'row') {
        return (
            <TableRow
                className="cursor-pointer hover:bg-muted/50"
                onClick={handleClick}
            >
                {showIcon && (
                    <TableCell>
                        <Box className="h-5 w-5 text-muted-foreground" />
                    </TableCell>
                )}
                <TableCell>
                    <div className="font-medium">{title}</div>
                    {showDescription && (
                        <div className="text-sm text-muted-foreground">{description}</div>
                    )}
                </TableCell>
                {visibleFields.map((fieldName) => {
                    const item = items.find((item) => item.label === fieldName);
                    return (
                        <TableCell key={fieldName}>
                            {item?.value || '-'}
                        </TableCell>
                    );
                })}
            </TableRow>
        );
    }

    const orderedItems = visibleFields
        .filter(fieldName => fieldName !== 'Repository_name' && fieldName !== 'Description')
        .map(fieldName => items.find(item => item.label === fieldName))
        .filter((item): item is Item => !!item && !!item.value && item.value !== '-' && item.value !== '');

    return (
        <div
            className="rounded-lg border bg-card text-card-foreground shadow-sm cursor-pointer hover:bg-muted/50"
            onClick={handleClick}
        >
            <div className="p-6 space-y-4">
                <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                        {showIcon && <Box className="h-5 w-5 text-muted-foreground" />}
                        <h3 className="font-semibold leading-none tracking-tight">{title}</h3>
                    </div>
                    <div className={`${chipColor} px-2.5 py-0.5 rounded-full text-xs font-medium`}>
                        {chipText}
                    </div>
                </div>
                {showDescription && (
                    <p className="text-sm text-muted-foreground">{description}</p>
                )}
                {orderedItems.length > 0 && (
                    <div className="grid grid-cols-2 gap-4">
                        {orderedItems.map((item: Item, index: number) => (
                            <div key={`${item.label}-${index}`} className="space-y-1">
                                <p className="text-sm font-medium leading-none">{item.label}</p>
                                <p className="text-sm text-muted-foreground">{item.value}</p>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
};

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
            const filteredItems = item.items?.filter((itemField: Item) => {
                const fieldConfig = activeSet?.fields.find(f => f.name === itemField.label);
                return fieldConfig?.isVisibleInRow;
            }) ?? [];

            return (
                <Card
                    key={item.id}
                    id={item.id}
                    title={item.label} // Assuming label is the title
                    description={item.value} // Assuming value is the description
                    chipColor="bg-blue-500" // Example chip color
                    chipText="Example" // Example chip text
                    items={filteredItems}
                    sharedLabels={sharedLabels}
                    onSelect={onSelect}
                    visibleFields={visibleFields}
                    viewMode={viewMode}
                    showDescription={showDescription}
                    showIcon={showIcon}
                />
            );
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
                const filteredItems = item.items?.filter((itemField: Item) => {
                    const fieldConfig = activeSet?.fields.find(f => f.name === itemField.label);
                    return fieldConfig?.isVisibleInRow;
                }) ?? [];


                return (
                    <Card
                        key={item.id}
                        id={item.id}
                        title={item.label} // Assuming label is the title
                        description={item.value} // Assuming value is the description
                        chipColor="bg-blue-500" // Example chip color
                        chipText="Example" // Example chip text
                        items={filteredItems}
                        sharedLabels={sharedLabels}
                        onSelect={onSelect}
                        visibleFields={visibleFields}
                        viewMode="row"
                        showDescription={true}
                        showIcon={showIcon}
                    />
                );
            })}
        </TableBody>
    </Table>
);

const arePropsEqual = (prevProps: CardComponentProps, nextProps: CardComponentProps) => {
    if (
        prevProps.id !== nextProps.id ||
        prevProps.title !== nextProps.title ||
        prevProps.description !== nextProps.description ||
        prevProps.chipColor !== nextProps.chipColor ||
        prevProps.chipText !== nextProps.chipText ||
        prevProps.viewMode !== nextProps.viewMode ||
        prevProps.showDescription !== nextProps.showDescription ||
        prevProps.showIcon !== nextProps.showIcon
    ) {
        return false;
    }

    if (prevProps.items.length !== nextProps.items.length) {
        return false;
    }

    const itemsChanged = prevProps.items.some((item, index) => {
        const nextItem = nextProps.items[index];
        return (
            item.label !== nextItem.label ||
            item.value !== nextItem.value ||
            item.isHeader !== nextItem.isHeader
        );
    });

    if (itemsChanged) {
        return false;
    }

    if (prevProps.visibleFields.length !== nextProps.visibleFields.length) {
        return false;
    }

    const visibleFieldsChanged = prevProps.visibleFields.some(
        (field, index) => field !== nextProps.visibleFields[index]
    );

    if (visibleFieldsChanged) {
        return false;
    }

    if (prevProps.sharedLabels.length !== nextProps.sharedLabels.length) {
        return false;
    }

    const sharedLabelsChanged = prevProps.sharedLabels.some(
        (label, index) => label !== nextProps.sharedLabels[index]
    );

    if (sharedLabelsChanged) {
        return false;
    }

    return true;
};

const Card = React.memo(CardComponent, arePropsEqual);

export { Card, GridView, TableView };