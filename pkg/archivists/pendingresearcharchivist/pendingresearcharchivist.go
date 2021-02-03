package pendingresearcharchivist

/*
The pending research archivist does the following: WHen initialized, it
immediately calculates the work surplus. The work surplus is equal to the
number of downstream researchers minus the number of work-items currently on
the pending queue. If this number is negative or zero, the service terminates.
If the number is positive, the one work item is added to the queue and the
process repeats at the work surplus calculation step. Keeping each instance
focused on only adding one work item at a time reduces the overhead of each
instance and allows the process of work-creation to more effectively scale
horizaontally.

*/
