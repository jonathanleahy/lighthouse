
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { ArrowUpDown, Filter, LayoutGrid, LayoutList, Search, ChevronUp, ChevronDown } from "lucide-react";

interface DashboardGuideProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

export default function DashboardGuide({ open, onOpenChange }: DashboardGuideProps) {
    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-4xl">
                <DialogHeader>
                    <DialogTitle>Dashboard Controls Guide</DialogTitle>
                </DialogHeader>

                <div className="grid gap-12 grid-cols-2">
                    <div className="space-y-8">
                        <div>
                            <div className="flex items-center gap-2 mb-3">
                                <Search className="h-5 w-5 text-muted-foreground" />
                                <h3 className="font-semibold text-base">Search & Filters</h3>
                            </div>
                            <ul className="space-y-2 text-sm">
                                <li><code className="px-2 py-0.5 bg-muted rounded-sm">active</code> - Exact match</li>
                                <li><code className="px-2 py-0.5 bg-muted rounded-sm">&gt;100</code> - Greater than</li>
                                <li><code className="px-2 py-0.5 bg-muted rounded-sm">&lt;50</code> - Less than</li>
                                <li><code className="px-2 py-0.5 bg-muted rounded-sm">%smith%</code> - Contains text</li>
                                <li><code className="px-2 py-0.5 bg-muted rounded-sm">ABC%</code> - Starts with</li>
                                <li><code className="px-2 py-0.5 bg-muted rounded-sm">%.com</code> - Ends with</li>
                            </ul>
                        </div>

                        <div>
                            <div className="flex items-center gap-2 mb-3">
                                <Filter className="h-5 w-5 text-muted-foreground" />
                                <h3 className="font-semibold text-base">Combining Filters</h3>
                            </div>
                            <ul className="space-y-2 text-sm">
                                <li><code className="px-2 py-0.5 bg-muted rounded-sm">active&amp;&amp;&gt;3</code> - AND condition</li>
                                <li><code className="px-2 py-0.5 bg-muted rounded-sm">pending||review</code> - OR condition</li>
                                <li><code className="px-2 py-0.5 bg-muted rounded-sm">!active</code> - NOT condition</li>
                                <li><code className="px-2 py-0.5 bg-muted rounded-sm">(A||B)&amp;&amp;C</code> - Complex combinations</li>
                            </ul>
                        </div>
                    </div>

                    <div className="space-y-8">
                        <div>
                            <div className="flex items-center gap-2 mb-3">
                                <ArrowUpDown className="h-5 w-5 text-muted-foreground" />
                                <h3 className="font-semibold text-base">Sorting & Order</h3>
                            </div>
                            <ul className="space-y-2 text-sm">
                                <li>• Click once: Sort ascending (A to Z, 1 to 9)</li>
                                <li>• Click twice: Sort descending (Z to A, 9 to 1)</li>
                                <li>• Click third time: Remove sorting</li>
                                <li className="flex items-center gap-2 pt-2">
                                    <ChevronUp className="h-4 w-4" />
                                    <ChevronDown className="h-4 w-4" />
                                    Use these arrows in edit mode to reorder fields
                                </li>
                            </ul>
                        </div>

                        <div>
                            <div className="flex items-center gap-2 mb-3">
                                <div className="flex gap-1">
                                    <LayoutGrid className="h-5 w-5 text-muted-foreground" />
                                    <LayoutList className="h-5 w-5 text-muted-foreground" />
                                </div>
                                <h3 className="font-semibold text-base">Display Options</h3>
                            </div>
                            <ul className="space-y-2 text-sm">
                                <li>• Card View: Grid layout with detailed information</li>
                                <li>• Table View: Traditional column format</li>
                                <li>• Toggle visibility using the view icons</li>
                                <li>• Save custom layouts for different use cases</li>
                            </ul>
                        </div>
                    </div>
                </div>
            </DialogContent>
        </Dialog>
    );
}