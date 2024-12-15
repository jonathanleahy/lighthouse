"use client"

import { useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { ArrowLeft, ChevronUp, ChevronDown, Edit, Save, Filter, PlusCircle, Trash2, LayoutGrid, LayoutList, HelpCircle } from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent, CardDescription, CardFooter } from "@/components/ui/card";
import { Table as UITable, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { useCustomFields } from '@/lib/customFieldsContext';
import { ScrollArea } from "@/components/ui/scroll-area";
import DashboardGuide from "@/components/DashboardGuide";

type SortOrder = 'asc' | 'desc' | 'none';
type DisplayMode = 'card' | 'row';

interface CustomField {
    id: string;
    name: string;
    isVisibleInCard: boolean;
    isVisibleInRow: boolean;
    filterEnabled?: boolean;
    filter?: string;
    sortOrder?: SortOrder;
}

interface CustomFieldSet {
    id: string;
    name: string;
    fields: CustomField[];
    displayMode: DisplayMode;
}

interface EditingState {
    isEditMode: boolean;
    editingName: string;
}

interface VisibilityToggleProps {
    isVisibleInCard: boolean;
    isVisibleInRow: boolean;
    isEditMode: boolean;
    onToggleCardVisibility: () => void;
    onToggleRowVisibility: () => void;
}

interface FieldSetGroup {
    name: string;
    sets: CustomFieldSet[];
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
    const [editingState, setEditingState] = useState<EditingState>({
        isEditMode: false,
        editingName: ''
    });
    const [newSetName, setNewSetName] = useState('');
    const [showGuide, setShowGuide] = useState(false);

    const getActiveSet = () => state.customFieldSets.find(set => set.id === state.activeSetId) || state.customFieldSets[0];

    const handleSetNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setEditingState(prev => ({
            ...prev,
            editingName: e.target.value
        }));
    };

    const handleSave = () => {
        const activeSet = getActiveSet();
        if (activeSet) {
            const updatedSet = {
                ...activeSet,
                name: editingState.editingName
            };
            dispatch({
                type: 'UPDATE_CUSTOM_FIELD_SET',
                payload: updatedSet
            });
            localStorage.setItem('customFieldSets', JSON.stringify(
                state.customFieldSets.map(set =>
                    set.id === activeSet.id ? updatedSet : set
                )
            ));
            localStorage.setItem('activeSetId', state.activeSetId);
        }
        setEditingState({ isEditMode: false, editingName: '' });
    };

    const handleCancel = useCallback(() => {
        setEditingState({ isEditMode: false, editingName: '' });
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

    const handleFilterToggle = (fieldId: string) => {
        const activeSet = getActiveSet();
        if (!activeSet) return;

        const field = activeSet.fields.find(f => f.id === fieldId);
        if (field) {
            updateField(fieldId, {
                filterEnabled: !field.filterEnabled,
                filter: field.filterEnabled ? '' : field.filter || ''
            });
        }
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
        <div className="container mx-auto p-4">
            <div className="flex items-center mb-6">
                <Button variant="ghost" onClick={() => router.push('/')} className="mr-4">
                    <ArrowLeft className="h-4 w-4 mr-2"/>
                    Back to Dashboard
                </Button>
                <h1 className="text-2xl font-bold">Customize Fields</h1>
                <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => setShowGuide(true)}
                    className="ml-auto"
                    title="View Controls Guide"
                >
                    <HelpCircle className="h-5 w-5"/>
                </Button>
            </div>

            <div className={`flex justify-between items-center mb-4 px-4 ${editingState.isEditMode ? 'opacity-50' : ''}`}>
                <ScrollArea className="w-[calc(100%-120px)]">
                    <div className="flex space-x-2">
                        {state.customFieldSets.map(set => (
                            <Button
                                key={set.id}
                                variant={state.activeSetId === set.id ? "default" : "outline"}
                                onClick={() => dispatch({type: 'SET_ACTIVE_SET_ID', payload: set.id})}
                                disabled={editingState.isEditMode}
                            >
                                {set.name}
                            </Button>
                        ))}
                    </div>
                </ScrollArea>
                <Dialog>
                    <DialogTrigger asChild>
                        <Button variant="outline" disabled={editingState.isEditMode}>
                            <PlusCircle className="h-4 w-4 mr-2"/>
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

            <Card>
                <CardHeader>
                    <div className="flex justify-between items-center">
                        <div>
                            {editingState.isEditMode ? (
                                <div className="mb-4">
                                    <Label htmlFor="setName">Set Name</Label>
                                    <Input
                                        value={editingState.editingName}
                                        onChange={handleSetNameChange}
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
                            {editingState.isEditMode ? (
                                <Select
                                    value={getActiveSet()?.displayMode}
                                    onValueChange={(value: DisplayMode) => {
                                        const activeSet = getActiveSet();
                                        if (activeSet) {
                                            dispatch({
                                                type: 'UPDATE_CUSTOM_FIELD_SET',
                                                payload: { ...activeSet, displayMode: value }
                                            });
                                        }
                                    }}
                                >
                                    <SelectTrigger className="w-[180px]">
                                        <SelectValue placeholder="Select display mode"/>
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="row">
                                            <span className="flex items-center">
                                                <LayoutList className="h-4 w-4 mr-2"/>
                                                Row View
                                            </span>
                                        </SelectItem>
                                        <SelectItem value="card">
                                            <span className="flex items-center">
                                                <LayoutGrid className="h-4 w-4 mr-2"/>
                                                Card View
                                            </span>
                                        </SelectItem>
                                    </SelectContent>
                                </Select>
                            ) : (
                                <div className="flex items-center space-x-2">
                                    {getActiveSet()?.displayMode === 'row' ? (
                                        <LayoutList className="h-6 w-6"/>
                                    ) : (
                                        <LayoutGrid className="h-6 w-6"/>
                                    )}
                                </div>
                            )}
                        </div>
                    </div>
                </CardHeader>
                <CardContent>
                    <UITable>



                        <TableHeader>
                            <TableRow>
                                {editingState.isEditMode && <TableHead className="w-[80px]">Order</TableHead>}
                                <TableHead className="w-[300px]">Field</TableHead>
                                <TableHead className="w-[100px]">Card View</TableHead>
                                <TableHead className="w-[100px]">Table View</TableHead>
                                <TableHead className="w-[400px]">Custom Filter</TableHead>
                        <TableHead className="w-[180px]">Sort</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {getActiveSet()?.fields.map((field, index) => (
                        <TableRow key={field.id}>
                            {editingState.isEditMode && (
                                <TableCell className="w-[80px]">
                                    <div className="flex items-center gap-2">
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            onClick={() => moveField(index, 'up')}
                                            disabled={index === 0}
                                        >
                                            <ChevronUp className="h-4 w-4"/>
                                        </Button>
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            onClick={() => moveField(index, 'down')}
                                            disabled={index === getActiveSet()?.fields.length - 1}
                                        >
                                            <ChevronDown className="h-4 w-4"/>
                                        </Button>
                                    </div>
                                </TableCell>
                            )}
                            <TableCell className="w-[300px]">{field.name}</TableCell>
                            <TableCell className="w-[100px]">
                                {editingState.isEditMode ? (
                                    <Button
                                        variant="ghost"
                                        size="icon"
                                        onClick={() => toggleFieldVisibility(field.id, 'card')}
                                        className="transition-all duration-200 relative"
                                        title={field.isVisibleInCard ? 'Hide in card view' : 'Show in card view'}
                                    >
                                        <LayoutGrid
                                            className={`h-4 w-4 ${field.isVisibleInCard ? 'text-foreground' : 'text-muted-foreground'}`}
                                        />
                                        {!field.isVisibleInCard && (
                                            <div className="absolute inset-0 flex items-center justify-center">
                                                <div className="w-5 h-px bg-muted-foreground rotate-45 transform origin-center"/>
                                            </div>
                                        )}
                                    </Button>
                                ) : (
                                    <div className="p-2 relative">
                                        <LayoutGrid
                                            className={`h-4 w-4 ${field.isVisibleInCard ? 'text-foreground' : 'text-muted-foreground'}`}
                                        />
                                        {!field.isVisibleInCard && (
                                            <div className="absolute inset-0 flex items-center justify-center">
                                                <div className="w-5 h-px bg-muted-foreground rotate-45 transform origin-center"/>
                                            </div>
                                        )}
                                    </div>
                                )}
                            </TableCell>
                            <TableCell className="w-[100px]">
                                {editingState.isEditMode ? (
                                    <Button
                                        variant="ghost"
                                        size="icon"
                                        onClick={() => toggleFieldVisibility(field.id, 'row')}
                                        className="transition-all duration-200 relative"
                                        title={field.isVisibleInRow ? 'Hide in table view' : 'Show in table view'}
                                    >
                                        <LayoutList
                                            className={`h-4 w-4 ${field.isVisibleInRow ? 'text-foreground' : 'text-muted-foreground'}`}
                                        />
                                        {!field.isVisibleInRow && (
                                            <div className="absolute inset-0 flex items-center justify-center">
                                                <div className="w-5 h-px bg-muted-foreground rotate-45 transform origin-center"/>
                                            </div>
                                        )}
                                    </Button>
                                ) : (
                                    <div className="p-2 relative">
                                        <LayoutList
                                            className={`h-4 w-4 ${field.isVisibleInRow ? 'text-foreground' : 'text-muted-foreground'}`}
                                        />
                                        {!field.isVisibleInRow && (
                                            <div className="absolute inset-0 flex items-center justify-center">
                                                <div className="w-5 h-px bg-muted-foreground rotate-45 transform origin-center"/>
                                            </div>
                                        )}
                                    </div>
                                )}
                            </TableCell>
                            <TableCell className="w-[400px]">
                                {editingState.isEditMode ? (
                                    <div className="flex items-center gap-2">
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            onClick={() => handleFilterToggle(field.id)}
                                        >
                                            <Filter
                                                className={`h-4 w-4 ${field.filterEnabled ? 'text-foreground' : 'text-muted-foreground'}`}/>
                                        </Button>
                                        <Input
                                            value={field.filter || ''}
                                            onChange={(e) => handleFilterChange(field.id, e.target.value)}
                                            placeholder="Enter filter"
                                            disabled={!field.filterEnabled}
                                            className="w-full"
                                        />
                                    </div>
                                ) : (
                                    field.filterEnabled ? field.filter : 'N/A'
                                )}
                            </TableCell>
                            <TableCell className="w-[180px]">
                                {editingState.isEditMode ? (
                                    <Select
                                        value={field.sortOrder || 'none'}
                                        onValueChange={(value: SortOrder) => handleSortChange(field.id, value)}
                                    >
                                        <SelectTrigger className="w-[170px]">
                                            <SelectValue placeholder="Select sort order"/>
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
    <CardFooter className="flex justify-between mt-4">
        {editingState.isEditMode ? (
            <>
                {state.customFieldSets.length > 1 && (
                    <Button
                        variant="destructive"
                        onClick={() => deleteCustomFieldSet(getActiveSet()?.id || '')}
                        disabled={state.activeSetId === '1'}
                    >
                        <Trash2 className="h-4 w-4 mr-2"/>
                        Delete Set
                    </Button>
                )}
                <div className="space-x-2 ml-auto">
                    <Button onClick={handleCancel} variant="outline">
                        Cancel
                    </Button>
                    <Button onClick={handleSave}>
                        <Save className="h-4 w-4 mr-2"/>
                        Save Changes
                    </Button>
                </div>
            </>
        ) : (
            <div className="w-full flex justify-end">
                <Button
                    onClick={() => {
                        const activeSet = getActiveSet();
                        setEditingState({
                            isEditMode: true,
                            editingName: activeSet?.name || ''
                        });
                    }}
                    variant="outline"
                >
                    <Edit className="h-4 w-4 mr-2"/>
                    Edit Fields
                </Button>
            </div>
        )}
    </CardFooter>
</Card>
    <DashboardGuide open={showGuide} onOpenChange={setShowGuide} />
</div>
);
}

