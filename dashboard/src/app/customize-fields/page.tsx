"use client"

import { useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { ArrowLeft, ChevronUp, ChevronDown, Edit, Save, PlusCircle, Trash2, LayoutGrid, LayoutList } from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent, CardDescription, CardFooter } from "@/components/ui/card";
import { Table as UITable, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { useCustomFields, CustomField, CustomFieldSet } from '@/lib/customFieldsContext';
import { ScrollArea } from "@/components/ui/scroll-area";

type SortOrder = 'asc' | 'desc' | 'none';
type DisplayMode = 'card' | 'row';


interface VisibilityToggleProps {
    isVisibleInCard: boolean;
    isVisibleInRow: boolean;
    isEditMode: boolean;
    onToggleCardVisibility: () => void;
    onToggleRowVisibility: () => void;
}

const VisibilityToggle: React.FC<VisibilityToggleProps> = ({
                                                               isVisibleInCard,
                                                               isVisibleInRow,
                                                               isEditMode,
                                                               onToggleCardVisibility,
                                                               onToggleRowVisibility
                                                           }) => {
    if (isEditMode) {
        return (
            <div className="flex items-center gap-2">
                <Button
                    variant="ghost"
                    size="icon"
                    onClick={onToggleCardVisibility}
                    className="transition-all duration-200 relative group"
                    title={isVisibleInCard ? 'Hide in card view' : 'Show in card view'}
                >
                    <LayoutGrid
                        className={`h-4 w-4 ${isVisibleInCard ? 'text-foreground' : 'text-muted-foreground'}`}
                    />
                    {!isVisibleInCard && (
                        <div className="absolute inset-0 flex items-center justify-center">
                            <div className="w-5 h-px bg-muted-foreground rotate-45 transform origin-center" />
                        </div>
                    )}
                </Button>

                <Button
                    variant="ghost"
                    size="icon"
                    onClick={onToggleRowVisibility}
                    className="transition-all duration-200 relative group"
                    title={isVisibleInRow ? 'Hide in table view' : 'Show in table view'}
                >
                    <LayoutList
                        className={`h-4 w-4 ${isVisibleInRow ? 'text-foreground' : 'text-muted-foreground'}`}
                    />
                    {!isVisibleInRow && (
                        <div className="absolute inset-0 flex items-center justify-center">
                            <div className="w-5 h-px bg-muted-foreground rotate-45 transform origin-center" />
                        </div>
                    )}
                </Button>
            </div>
        );
    }

    return (
        <div className="flex items-center gap-2">
            <div className="p-2 relative">
                <LayoutGrid
                    className={`h-4 w-4 ${isVisibleInCard ? 'text-foreground' : 'text-muted-foreground'}`}
                />
                {!isVisibleInCard && (
                    <div className="absolute inset-0 flex items-center justify-center">
                        <div className="w-5 h-px bg-muted-foreground rotate-45 transform origin-center" />
                    </div>
                )}
            </div>

            <div className="p-2 relative">
                <LayoutList
                    className={`h-4 w-4 ${isVisibleInRow ? 'text-foreground' : 'text-muted-foreground'}`}
                />
                {!isVisibleInRow && (
                    <div className="absolute inset-0 flex items-center justify-center">
                        <div className="w-5 h-px bg-muted-foreground rotate-45 transform origin-center" />
                    </div>
                )}
            </div>
        </div>
    );
};

export default function CustomizeFieldsPage() {
    const router = useRouter();
    const { state, dispatch } = useCustomFields();
    const [isEditMode, setIsEditMode] = useState(false);
    const [newSetName, setNewSetName] = useState('');

    const getActiveSet = () => state.customFieldSets.find(set => set.id === state.activeSetId) || state.customFieldSets[0];

    const handleSave = () => {
        const activeSet = getActiveSet();
        if (activeSet) {
            localStorage.setItem('customFieldSets', JSON.stringify(state.customFieldSets));
            localStorage.setItem('activeSetId', state.activeSetId);
        }
        setIsEditMode(false);
    };

    const handleCancel = useCallback(() => {
        setIsEditMode(false);
    }, []);

    const moveField = (index: number, direction: 'up' | 'down') => {
        const activeSet = getActiveSet();
        if (!activeSet) return;

        const newFields = [...activeSet.fields];
        const newIndex = direction === 'up' ? index - 1 : index + 1;
        if (newIndex >= 0 && newIndex < newFields.length) {
            [newFields[index], newFields[newIndex]] = [newFields[newIndex], newFields[index]];
            dispatch({
                type: 'UPDATE_CUSTOM_FIELD_SET',
                payload: { ...activeSet, fields: newFields }
            });
        }
    };

    const handleSortChange = (fieldId: string, order: SortOrder) => {
        updateField(fieldId, { sortOrder: order });
    };

    const handleFilterChange = (fieldId: string, value: string) => {
        updateField(fieldId, { filter: value });
    };

    const toggleFieldVisibility = (fieldId: string, viewMode: 'card' | 'row') => {
        const activeSet = getActiveSet();
        if (!activeSet) return;

        const field = activeSet.fields.find(f => f.id === fieldId);
        if (field) {
            const updates = viewMode === 'card'
                ? { isVisibleInCard: !field.isVisibleInCard }
                : { isVisibleInRow: !field.isVisibleInRow };

            updateField(fieldId, updates);
        }
    };

    const updateField = (fieldId: string, updates: Partial<CustomField>) => {
        const activeSet = getActiveSet();
        if (!activeSet) return;

        const updatedFields = activeSet.fields.map(field =>
            field.id === fieldId ? { ...field, ...updates } : field
        );
        dispatch({
            type: 'UPDATE_CUSTOM_FIELD_SET',
            payload: { ...activeSet, fields: updatedFields }
        });
    };

    const addCustomFieldSet = () => {
        if (newSetName.trim()) {
            const newSet: CustomFieldSet = {
                id: Date.now().toString(),
                name: newSetName.trim(),
                fields: getActiveSet()?.fields.map(field => ({ ...field, id: `${Date.now()}-${field.id}` })) || [],
                displayMode: 'row',
                viewConfig: {
                    defaultView: 'table',
                    cardConfig: {
                        showDescription: true,
                        showIcon: true
                    },
                    tableConfig: {
                        compact: false,
                        showIcon: true
                    }
                }
            };
            dispatch({ type: 'ADD_CUSTOM_FIELD_SET', payload: newSet });
            dispatch({ type: 'SET_ACTIVE_SET_ID', payload: newSet.id });
            setNewSetName('');
        }
    };

    const deleteCustomFieldSet = (id: string) => {
        if (state.customFieldSets.length > 1) {
            dispatch({ type: 'DELETE_CUSTOM_FIELD_SET', payload: id });
        }
    };

    return (
        <div className="container mx-auto p-4 h-screen flex flex-col">
            <div className="flex items-center mb-6">
                <Button variant="ghost" onClick={() => router.push('/')} className="mr-4">
                    <ArrowLeft className="h-4 w-4 mr-2" />
                    Back to Dashboard
                </Button>
                <h1 className="text-2xl font-bold">Customize Fields</h1>
            </div>

            <div className={`flex justify-between items-center mb-4 px-4 ${isEditMode ? 'opacity-50' : ''}`}>
                <ScrollArea className="w-[calc(100%-120px)]">
                    <div className="flex space-x-2">
                        {state.customFieldSets.map(set => (
                            <Button
                                key={set.id}
                                variant={state.activeSetId === set.id ? "default" : "outline"}
                                onClick={() => dispatch({ type: 'SET_ACTIVE_SET_ID', payload: set.id })}
                                disabled={isEditMode}
                            >
                                {set.name}
                            </Button>
                        ))}
                    </div>
                </ScrollArea>
                <Dialog>
                    <DialogTrigger asChild>
                        <Button variant="outline" disabled={isEditMode}>
                            <PlusCircle className="h-4 w-4 mr-2" />
                            Add New Set
                        </Button>
                    </DialogTrigger>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>Add New Field Set</DialogTitle>
                        </DialogHeader>
                        <div className="flex items-center space-x-2">
                            <Input
                                value={newSetName}
                                onChange={(e) => setNewSetName(e.target.value)}
                                placeholder="Enter set name"
                            />
                            <Button onClick={addCustomFieldSet}>Add</Button>
                        </div>
                    </DialogContent>
                </Dialog>
            </div>

            <Card className="flex-grow flex flex-col">
                <CardHeader>
                    <div className="flex justify-between items-center">
                        <div>
                            {isEditMode ? (
                                <div className="mb-4">
                                    <Label htmlFor="setName">Set Name</Label>
                                    <Input
                                        value={getActiveSet()?.name || ''}
                                        onChange={(e) => updateField(getActiveSet()?.id || '', { name: e.target.value })}
                                        id="setName"
                                    />
                                </div>
                            ) : (
                                <>
                                    <CardTitle>{getActiveSet()?.name}</CardTitle>
                                    <CardDescription>Customize the fields for this set</CardDescription>
                                </>
                            )}
                        </div>
                        <div className="flex items-center space-x-2">
                            {isEditMode ? (
                                <Select
                                    value={getActiveSet()?.displayMode}
                                    onValueChange={(value: DisplayMode) => {
                                        updateField(getActiveSet()?.id || '', { displayMode: value });
                                    }}
                                >
                                    <SelectTrigger className="w-[180px]">
                                        <SelectValue placeholder="Select display mode" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="row">
                                            Row
                                        </SelectItem>
                                        <SelectItem value="card">
                                            Card
                                        </SelectItem>
                                    </SelectContent>
                                </Select>
                            ) : (
                                <div className="flex items-center space-x-2">
                                    {getActiveSet()?.displayMode === 'row' ? (
                                        <LayoutList className="h-6 w-6" />
                                    ) : (
                                        <LayoutGrid className="h-6 w-6" />
                                    )}
                                </div>
                            )}
                        </div>
                    </div>
                </CardHeader>
                <CardContent className="flex-grow overflow-auto">
                    {!isEditMode && (
                        <div className="mb-4">
                            <Button onClick={() => setIsEditMode(true)} variant="outline">
                                <Edit className="h-4 w-4 mr-2" />
                                Edit Fields
                            </Button>
                        </div>
                    )}
                    <UITable>
                        <TableHeader>
                            <TableRow>
                                {isEditMode && <TableHead className="w-[100px]">Order</TableHead>}
                                <TableHead className="w-[200px]">Field</TableHead>
                                <TableHead className="w-[100px]">Visible</TableHead>
                                <TableHead className="w-[200px]">Custom Filter</TableHead>
                                <TableHead className="w-[100px]">Sort</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {getActiveSet()?.fields.map((field, index) => (
                                <TableRow key={field.id}>
                                    {isEditMode && (
                                        <TableCell>
                                            <div className="flex items-center gap-2">
                                                <Button
                                                    variant="ghost"
                                                    size="icon"
                                                    onClick={() => moveField(index, 'up')}
                                                    disabled={index === 0}
                                                >
                                                    <ChevronUp className="h-4 w-4" />
                                                </Button>
                                                <Button
                                                    variant="ghost"
                                                    size="icon"
                                                    onClick={() => moveField(index, 'down')}
                                                    disabled={index === getActiveSet()?.fields.length - 1}
                                                >
                                                    <ChevronDown className="h-4 w-4" />
                                                </Button>
                                            </div>
                                        </TableCell>
                                    )}
                                    <TableCell>{field.name}</TableCell>
                                    <TableCell>
                                        <VisibilityToggle
                                            isVisibleInCard={field.isVisibleInCard}
                                            isVisibleInRow={field.isVisibleInRow}
                                            isEditMode={isEditMode}
                                            onToggleCardVisibility={() => toggleFieldVisibility(field.id, 'card')}
                                            onToggleRowVisibility={() => toggleFieldVisibility(field.id, 'row')}
                                        />
                                    </TableCell>
                                    <TableCell>
                                        {isEditMode ? (
                                            <Input
                                                value={field.filter || ''}
                                                onChange={(e) => handleFilterChange(field.id, e.target.value)}
                                                placeholder="Enter filter"
                                            />
                                        ) : (
                                            field.filter
                                        )}
                                    </TableCell>
                                    <TableCell>
                                        {isEditMode ? (
                                            <Select
                                                value={field.sortOrder || 'none'}
                                                onValueChange={(value: SortOrder) => handleSortChange(field.id, value)}
                                            >
                                                <SelectTrigger className="w-[100px]">
                                                    <SelectValue placeholder="Select sort order" />
                                                </SelectTrigger>
                                                <SelectContent>
                                                    <SelectItem value="none">None</SelectItem>
                                                    <SelectItem value="asc">Ascending</SelectItem>
                                                    <SelectItem value="desc">Descending</SelectItem>
                                                </SelectContent>
                                            </Select>
                                        ) : (
                                            field.sortOrder
                                        )}
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </UITable>
                </CardContent>
                {isEditMode && (
                    <CardFooter className="flex justify-between">
                        {state.customFieldSets.length > 1 && (
                            <Button
                                variant="destructive"
                                onClick={() => deleteCustomFieldSet(getActiveSet()?.id || '')}
                                disabled={state.activeSetId === '1'}
                            >
                                <Trash2 className="h-4 w-4 mr-2" />
                                Delete Set
                            </Button>
                        )}
                        <div className="space-x-2 ml-auto">
                            <Button onClick={handleCancel} variant="outline">
                                Cancel
                            </Button>
                            <Button onClick={handleSave}>
                                <Save className="h-4 w-4 mr-2" />
                                Save Changes
                            </Button>
                        </div>
                    </CardFooter>
                )}
            </Card>
        </div>
    );
}