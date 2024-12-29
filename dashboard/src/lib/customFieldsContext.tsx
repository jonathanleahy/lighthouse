"use client"

import {createContext, useContext, useReducer, ReactNode, useState, useEffect} from 'react'

export type SortOrder = 'asc' | 'desc' | 'none'
export type DisplayMode = 'card' | 'row'

export type CustomField = {
  id: string;
  name: string;
  isVisibleInCard: boolean;
  isVisibleInRow: boolean;
  sortOrder: 'asc' | 'desc' | 'none';
  filter: string;
  filterEnabled: boolean;
  type: string;
  displayMode: 'card' | 'row';
}

export type CustomFieldSet = {
  id: string
  name: string
  fields: CustomField[]
  displayMode: DisplayMode
  viewConfig: ViewConfig;
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

type CustomFieldsState = {
  customFieldSets: CustomFieldSet[]
  activeSetId: string
}

export type CustomFieldsAction =
    | { type: 'UPDATE_CUSTOM_FIELD_SET'; payload: CustomFieldSet }
    | { type: 'SET_ACTIVE_SET_ID'; payload: string }
    | { type: 'ADD_CUSTOM_FIELD_SET'; payload: CustomFieldSet }
    | { type: 'DELETE_CUSTOM_FIELD_SET'; payload: string }
    | { type: 'SET_CUSTOM_FIELD_SETS'; payload: CustomFieldSet[] }

function customFieldsReducer(state: CustomFieldsState, action: CustomFieldsAction): CustomFieldsState {
  switch (action.type) {
    case 'UPDATE_CUSTOM_FIELD_SET':
      return {
        ...state,
        customFieldSets: state.customFieldSets.map(set =>
            set.id === action.payload.id ? action.payload : set
        )
      }
    case 'SET_ACTIVE_SET_ID':
      return {
        ...state,
        activeSetId: action.payload
      }
    case 'ADD_CUSTOM_FIELD_SET':
      return {
        ...state,
        customFieldSets: [...state.customFieldSets, action.payload]
      }
    case 'DELETE_CUSTOM_FIELD_SET':
      const updatedSets = state.customFieldSets.filter(set => set.id !== action.payload)
      const newActiveId = state.activeSetId === action.payload ? updatedSets[0]?.id : state.activeSetId
      return {
        customFieldSets: updatedSets,
        activeSetId: newActiveId
      }
    case 'SET_CUSTOM_FIELD_SETS':
      return {
        ...state,
        customFieldSets: action.payload
      }
    default:
      return state
  }
}

const CustomFieldsContext = createContext<{
  state: CustomFieldsState
  dispatch: React.Dispatch<CustomFieldsAction>
} | null>(null)

export function CustomFieldsProvider({ children }: { children: ReactNode }) {
  // Use useState for initial loading state
  const [isLoading, setIsLoading] = useState(true)

  const [state, dispatch] = useReducer(customFieldsReducer, {
    customFieldSets: [],
    activeSetId: ''
  })

  useEffect(() => {
    // Try to load from localStorage first
    const savedSets = localStorage.getItem('customFieldSets')
    const savedActiveId = localStorage.getItem('activeSetId')

    if (savedSets && savedActiveId) {
      dispatch({ type: 'SET_CUSTOM_FIELD_SETS', payload: JSON.parse(savedSets) })
      dispatch({ type: 'SET_ACTIVE_SET_ID', payload: savedActiveId })
    }

    setIsLoading(false)
  }, [])

  if (isLoading) {
    return null
  }

  return (
      <CustomFieldsContext.Provider value={{ state, dispatch }}>
        {children}
      </CustomFieldsContext.Provider>
  )
}

export function useCustomFields() {
  const context = useContext(CustomFieldsContext)
  if (!context) {
    throw new Error('useCustomFields must be used within a CustomFieldsProvider')
  }
  return context
}