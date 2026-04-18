// Package diff provides utilities for processing, rendering, and exporting
// infrastructure drift results produced by the detector.
//
// Key components:
//   - Renderer   – unified-diff-like human-readable output
//   - Highlighter – colorises individual attribute changes
//   - Pager       – paginates large result sets
//   - Searcher    – filters results by type, ID, or kind
//   - Sorter      – sorts results by various fields
//   - Grouper     – groups results by resource type or change kind
//   - Annotator   – attaches explanatory notes to changes
//   - Comparator  – compares two snapshots of drift results over time
//   - Merger      – deduplicates and merges overlapping result sets
//   - Truncator   – limits output size for large environments
//   - Exporter    – serialises results to text or JSON
//   - Stats       – computes aggregate statistics over a result set
package diff
