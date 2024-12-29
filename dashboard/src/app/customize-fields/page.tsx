"use client"

import React, { useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import {
    ArrowLeft,
    ChevronUp,
    ChevronDown,
    Edit,
    Filter,
    PlusCircle,
    Trash2,
    LayoutGrid,
    HelpCircle,
    Check,
    LayoutList
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent, CardDescription, CardFooter } from "@/components/ui/card";
import { Table as UITable, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { useCustomFields } from '@/lib/customFieldsContext';
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
    type?: string; // Add missing properties
    displayMode?: DisplayMode;
}

import { CustomFieldSet as ImportedCustomFieldSet } from '@/lib/customFieldsContext';
import IconWithSlash from "@/components/IconWithSlash";
import Link from "next/link";
import {ThemeToggle} from "@/components/theme-toggle";

type CustomFieldSet = ImportedCustomFieldSet

interface EditingState {
    isEditMode: boolean;
    editingName: string;
}

export default function CustomizeFieldsPage() {
    const router = useRouter();
    const {state, dispatch} = useCustomFields();
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
        setEditingState({isEditMode: false, editingName: ''});
    };

    const handleCancel = useCallback(() => {
        setEditingState({isEditMode: false, editingName: ''});
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
                payload: {...activeSet, fields: newFields}
            });
        }
    };

    const handleSortChange = (fieldId: string, order: SortOrder) => {
        updateField(fieldId, {sortOrder: order});
    };

    const handleFilterChange = (fieldId: string, value: string) => {
        updateField(fieldId, {filter: value});
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
                ? {isVisibleInCard: !field.isVisibleInCard}
                : {isVisibleInRow: !field.isVisibleInRow};

            updateField(fieldId, updates);
        }
    };

    const updateField = (fieldId: string, updates: Partial<CustomField>) => {
        const activeSet = getActiveSet();
        if (!activeSet) return;

        const updatedFields = activeSet.fields.map(field =>
            field.id === fieldId ? {...field, ...updates} : field
        );
        dispatch({
            type: 'UPDATE_CUSTOM_FIELD_SET',
            payload: {...activeSet, fields: updatedFields}
        });
    };

    const addCustomFieldSet = () => {
        if (newSetName.trim()) {
            const newSet: CustomFieldSet = {
                id: Date.now().toString(),
                name: newSetName.trim(),
                fields: getActiveSet()?.fields.map(field => ({...field, id: `${Date.now()}-${field.id}`})) || [],
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
            dispatch({type: 'ADD_CUSTOM_FIELD_SET', payload: newSet});
            dispatch({type: 'SET_ACTIVE_SET_ID', payload: newSet.id});
            setNewSetName('');
        }
    };

    const deleteCustomFieldSet = (id: string) => {
        if (state.customFieldSets.length > 1) {
            dispatch({type: 'DELETE_CUSTOM_FIELD_SET', payload: id});
        }
    };

    return (
<div>
        {/*<div className="min-h-screen bg-background mb-6">*/}
    <header
        className="sticky top-0 z-20 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container mx-auto px-4 sm:px-4 lg:px-4 flex h-12 items-center justify-between">
            <div className="flex items-center space-x-4">
                <Link href="/" className="flex items-center space-x-2">
                    <ArrowLeft className="h-5 w-5"/>
                    <span className="font-medium">Dashboard</span>
                </Link>
                <span className="text-muted-foreground">/</span>
                <h1 className="text-xl font-bold">Custom Search Groups</h1>
            </div>
            <div className="flex items-center space-x-4">
                <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => setShowGuide(true)}
                    className="ml-auto"
                    title="View Controls Guide"
                >
                    <HelpCircle className="h-5 w-5"/>
                </Button>
                <ThemeToggle/>
            </div>
        </div>
    </header>

    <div className="container mx-auto p-4">

        {/*<div className="flex items-center mb-6">*/}
        {/*    <Button variant="ghost" onClick={() => router.push('/')} className="mr-4">*/}
        {/*        <ArrowLeft className="h-4 w-4 mr-2"/>*/}
        {/*        Back to Dashboard*/}
        {/*    </Button>*/}
        {/*    <h1 className="text-2xl font-bold">Customize Fields</h1>*/}
        {/*    <Button*/}
        {/*        variant="ghost"*/}
        {/*        size="icon"*/}
        {/*        onClick={() => setShowGuide(true)}*/}
        {/*        className="ml-auto"*/}
        {/*        title="View Controls Guide"*/}
        {/*    >*/}
        {/*        <HelpCircle className="h-5 w-5"/>*/}
        {/*    </Button>*/}
        {/*</div>*/}

        <div
            className={`flex justify-between items-center mb-4 px-4 ${editingState.isEditMode ? 'opacity-50' : ''}`}>
            <div className="w-full overflow-x-auto">
                <div className="inline-flex rounded-lg bg-muted p-1 text-muted-foreground">
                    {state.customFieldSets.map(set => (
                        <Button
                            key={set.id}
                            variant="ghost"
                            className={`rounded-md px-3 py-1.5 text-sm font-medium transition-all ${
                                state.activeSetId === set.id
                                    ? "bg-background text-foreground shadow-sm"
                                    : "hover:bg-background/50 hover:text-foreground"
                            }`}
                            onClick={() => dispatch({type: 'SET_ACTIVE_SET_ID', payload: set.id})}
                            disabled={editingState.isEditMode}
                        >
                            {set.name}
                        </Button>
                    ))}
                </div>
            </div>
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
                                <CardDescription>Customize the fields for this search set</CardDescription>
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
                                            payload: {...activeSet, displayMode: value}
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
                        {editingState.isEditMode ? (
                            <>
                                <Button variant="outline" onClick={handleCancel}>
                                    <ArrowLeft className="h-4 w-4 mr-2"/>
                                    Cancel
                                </Button>
                                <Button onClick={handleSave}>
                                    <Check className="h-4 w-4 mr-2"/>
                                    Save
                                </Button>
                            </>
                        ) : (
                            <Button variant="outline" onClick={() => setEditingState({
                                isEditMode: true,
                                editingName: getActiveSet()?.name || ''
                            })}>
                                <Edit className="h-4 w-4 mr-2"/>
                                Edit
                            </Button>
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
                                            onClick={() => toggleFieldVisibility(field.id, 'row')}
                                            className="transition-all duration-200"
                                            title={field.isVisibleInCard ? 'Hide in table view' : 'Show in table view'}
                                        >
                                            <IconWithSlash Icon={LayoutGrid} disabled={!field.isVisibleInRow}/>
                                        </Button>
                                    ) : (
                                        <IconWithSlash Icon={LayoutGrid} disabled={!field.isVisibleInRow}/>
                                    )}
                                </TableCell>
                                <TableCell className="w-[100px]">
                                    {editingState.isEditMode ? (
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            onClick={() => toggleFieldVisibility(field.id, 'card')}
                                            className="transition-all duration-200"
                                            title={field.isVisibleInCard ? 'Hide in card view' : 'Show in card view'}
                                        >
                                            <IconWithSlash Icon={LayoutGrid} disabled={!field.isVisibleInCard}/>
                                        </Button>
                                    ) : (
                                        <IconWithSlash Icon={LayoutGrid} disabled={!field.isVisibleInCard}/>
                                    )}
                                </TableCell>
                                <TableCell className="w-[400px]">
                                    {editingState.isEditMode ? (
                                        <div className="flex items-center gap-2">
                                            <Button
                                                variant="ghost"
                                                size="icon"
                                                onClick={() => handleFilterToggle(field.id)}
                                                className="transition-all duration-200 relative"
                                                title={field.filterEnabled ? 'Disable filter' : 'Enable filter'}
                                            >
                                                <Filter
                                                    className={`h-4 w-4 ${field.filterEnabled ? 'text-foreground' : 'text-muted-foreground'}`}
                                                />
                                                {!field.filterEnabled && (
                                                    <div
                                                        className="absolute inset-0 flex items-center justify-center">
                                                        <div
                                                            className="w-5 h-px bg-muted-foreground rotate-45 transform origin-center"/>
                                                    </div>
                                                )}
                                            </Button>
                                            <Input
                                                value={field.filter || ''}
                                                onChange={(e) => handleFilterChange(field.id, e.target.value)}
                                                placeholder="Enter filter value"
                                                disabled={!field.filterEnabled}
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
                                            <SelectTrigger className="w-full">
                                                <SelectValue placeholder="Select sort order"/>
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem value="none">None</SelectItem>
                                                <SelectItem value="asc">Ascending</SelectItem>
                                                <SelectItem value="desc">Descending</SelectItem>
                                            </SelectContent>
                                        </Select>
                                    ) : (
                                        field.sortOrder || 'N/A'
                                    )}
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </UITable>
            </CardContent>
            <CardFooter className="flex justify-between">
                <Button variant="destructive" onClick={() => deleteCustomFieldSet(getActiveSet()?.id || '')}>
                    <Trash2 className="h-4 w-4 mr-2"/>
                    Delete
                </Button>
            </CardFooter>
        </Card>

        {showGuide && <DashboardGuide open={showGuide} onOpenChange={setShowGuide}/>}
    </div>
</div>
)}

